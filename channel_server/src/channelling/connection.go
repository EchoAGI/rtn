package channelling

import (
	"container/list"
	"io"
	"log"
	"sync"
	"time"

	"buffercache"

	"github.com/gorilla/websocket"
)

const (
	// 写超时 to client.
	writeWait = 10 * time.Second

	// 读心跳信息的超时.
	pongWait = 60 * time.Second

	// 发送 pings 到 client 的超时. 必须小于 readWait.
	pingPeriod = (pongWait * 9) / 10

	// client 允许的最大消息数.
	maxMessageSize = 1024 * 1024

	// 发送队列的大小.
	queueSize    = 512
	maxQueueSize = queueSize * 4

	// Throttle 流量控制.
	maxRatePerSecond = 20
)



type Connection interface {
	Index() uint64
	Send(*Message)
	Close()
	ReadPump()
	WritePump()
}

type ConnectionHandler interface {
	NewBuffer() buffercache.Buffer
	OnConnect(Connection)
	OnDisconnect()
	OnText(buffercache.Buffer)
}

type connection struct {
	// References.
	ws      *websocket.Conn
	handler ConnectionHandler

	// Data handling.
	condition *sync.Cond
	queue     list.List
	mutex     sync.Mutex
	isClosed  bool

	// Debugging
	Idx uint64
}

func NewConnection(index uint64, ws *websocket.Conn, handler ConnectionHandler) Connection {
	c := &connection{
		ws:      ws,
		handler: handler,
		Idx:     index,
	}
	c.condition = sync.NewCond(&c.mutex)

	return c
}

func (c *connection) Index() uint64 {
	return c.Idx
}

func (c *connection) Close() {
	c.mutex.Lock()
	if c.isClosed {
		c.mutex.Unlock()
		return
	}
	
	c.isClosed = true
	c.mutex.Unlock()
	// Unlock while we close the websocket connection.
	c.ws.Close()


	// 再次锁定, 清理队列, 发送signal.
	c.mutex.Lock()
	for {
		head := c.queue.Front()
		if head == nil {
			break
		}
		c.queue.Remove(head)
		message := head.Value.(buffercache.Buffer)
		message.Decref()
	}
	c.condition.Signal()
	c.mutex.Unlock()
}

// readPump 从 websocket 读出消息到hub.
func (c *connection) ReadPump() {
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	times := list.New()

	// NOTE(lcooper): This more or less assumes that the write pump is started.
	c.handler.OnConnect(c)

	for {
		//fmt.Println("readPump wait nextReader", c.Idx)
		op, r, err := c.ws.NextReader()
		if err != nil {
			if err == io.EOF {
			} else {
				log.Println("Error while reading", c.Idx, err)
			}
			break
		}

		switch op {
		case websocket.TextMessage:
			now := time.Now()
			if times.Len() == maxRatePerSecond {
				front := times.Front()
				times.Remove(front)
				delta := time.Second - now.Sub(front.Value.(time.Time))
				if delta > 0 {
					// client is sending messages too fast, delay him
					time.Sleep(delta)
				}
			}
			times.PushBack(now)

			message := c.handler.NewBuffer()
			err = buffercache.ReadAll(message, r)
			if err != nil {
				message.Decref()
				break
			}
			c.handler.OnText(message)
			message.Decref()
		}
	}

	c.Close()
	c.handler.OnDisconnect()
}

// 把消息放入写出队列.
func (c *connection) Send(message *Message) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.isClosed {
		return
	}
	//fmt.Println("Outbound queue size", c.Idx, len(c.queue))
	if c.queue.Len() >= maxQueueSize {
		log.Println("Outbound queue overflow", c.Idx, c.queue.Len())
		return
	}
	message.Incref()
	c.queue.PushBack(message)
	c.condition.Signal()
}

// writePump pumps messages from the queue to the websocket connection.
func (c *connection) WritePump() {
	var timer *time.Timer
	ping := false

	// Spawn a timer to emit websocket pings.
	timer = time.AfterFunc(pingPeriod, func() {
		c.mutex.Lock()
		if c.isClosed {
			c.mutex.Unlock()
			return
		}
		ping = true
		c.condition.Signal()
		c.mutex.Unlock()
		timer.Reset(pingPeriod)
	})

	// Wait for actions.
	for {

		c.mutex.Lock()
		// Wait until something todo.
		for !ping && !c.isClosed && c.queue.Len() == 0 {
			// Wait on signal (this also unlocks while waiting, and locks again when got the signal).
			c.condition.Wait()
		}
		// Fast exit if in closed state.
		if c.isClosed {
			c.mutex.Unlock()
			goto cleanup
		}
		// Flush queue if something.
		for {
			head := c.queue.Front()
			if head == nil {
				break
			}
			c.queue.Remove(head)
			message := head.Value.(*Message)
			if ping {
				// Send ping.
				ping = false
				c.mutex.Unlock()
				if err := c.ping(); err != nil {
					log.Println("Error while sending ping", c.Idx, err)
					message.Decref()
					goto cleanup
				}
			} else {
				c.mutex.Unlock()
			}

			//websocket.BinaryMessage
			if err := c.write(message.FrameType, message.Bytes()); err != nil {
				log.Println("Error while writing", c.Idx, err)
				message.Decref()
				goto cleanup
			}
			message.Decref()
			c.mutex.Lock()
		}
		if ping {
			// Send ping.
			ping = false
			c.mutex.Unlock()
			if err := c.ping(); err != nil {
				log.Println("Error while sending ping", c.Idx, err)
				goto cleanup
			}
		} else {
			// Final unlock.
			c.mutex.Unlock()
		}

	}

cleanup:
	//fmt.Println("writePump done")
	timer.Stop()
	c.Close()
}

// Write ping message.
func (c *connection) ping() error {
	return c.write(websocket.PingMessage, []byte{})
}

// Write writes a message with the given opCode and payload.
func (c *connection) write(opCode int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(opCode, payload)
}
