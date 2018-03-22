package channelling

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats"

	"natsconnection"
)

const (
	BusManagerStartup    = "startup"
	BusManagerOffer      = "offer"
	BusManagerAnswer     = "answer"
	BusManagerBye        = "bye"
	BusManagerConnect    = "connect"
	BusManagerDisconnect = "disconnect"
	BusManagerSession    = "session"
)

// BusManager 提供了与 消息总线进行通信的API.
type BusManager interface {
	ChannellingAPIConsumer
	Start()
	Publish(subject string, v interface{}) error
	Request(subject string, v interface{}, vPtr interface{}, timeout time.Duration) error
	Trigger(name, from, payload string, data interface{}, pipeline *Pipeline) error
	Subscribe(subject string, cb nats.Handler) (*nats.Subscription, error)
	BindRecvChan(subject string, channel interface{}) (*nats.Subscription, error)
	BindSendChan(subject string, channel interface{}) error
	PrefixSubject(string) string
	CreateSink(string) Sink
}

// BusTrigger 作为序列化 后端系统总线 trigger 事件的容器.
type BusTrigger struct {
	Id       string
	Name     string
	From     string
	Payload  string      `json:",omitempty"`
	Data     interface{} `json:",omitempty"`
	Pipeline string      `json:",omitempty"`
}

// BusSubjectTrigger 返回 trigger payloads 的消息主题名称.
func BusSubjectTrigger(prefix, suffix string) string {
	return fmt.Sprintf("%s.%s", prefix, suffix)
}

// NewBusManager 创建和初始化一个新的 BusManager, 根据 useNats开关决定是否使用 NATS.
// 目的是为了简化API, 封装与后端消息总线进行连接和收发数据的逻辑.
func NewBusManager(apiConsumer ChannellingAPIConsumer, id string, useNats bool, subjectPrefix string) BusManager {
	var b BusManager
	var err error
	if useNats {
		b, err = newNatsBus(apiConsumer, id, subjectPrefix)
		if err == nil {
			log.Println("NATS bus connected")
		} else {
			log.Println("Error connecting NATS bus", err)
			b = &noopBus{apiConsumer, id}
		}
	} else {
		b = &noopBus{apiConsumer, id}
	}

	return b
}

type noopBus struct {
	ChannellingAPIConsumer
	id string
}

func (bus *noopBus) Start() {
	// noop
}

func (bus *noopBus) Publish(subject string, v interface{}) error {
	return nil
}

func (bus *noopBus) Request(subject string, v interface{}, vPtr interface{}, timeout time.Duration) error {
	return nil
}

func (bus *noopBus) Trigger(name, from, payload string, data interface{}, pipeline *Pipeline) error {
	return nil
}

func (bus *noopBus) PrefixSubject(subject string) string {
	return subject
}

func (bus *noopBus) BindRecvChan(subject string, channel interface{}) (*nats.Subscription, error) {
	return nil, nil
}

func (bus *noopBus) BindSendChan(subject string, channel interface{}) error {
	return nil
}

func (bus *noopBus) Subscribe(subject string, cb nats.Handler) (*nats.Subscription, error) {
	return nil, nil
}

func (bus *noopBus) CreateSink(id string) Sink {
	return nil
}

type natsBus struct {
	ChannellingAPIConsumer
	id           string
	prefix       string
	ec           *natsconnection.EncodedConnection
	triggerQueue chan *busQueueEntry
}

func newNatsBus(apiConsumer ChannellingAPIConsumer, id, prefix string) (*natsBus, error) {
	ec, err := natsconnection.EstablishJSONEncodedConnection(nil)
	if err != nil {
		return nil, err
	}
	if prefix == "" {
		prefix = "channelling.trigger"
	}
	// Create buffered channel for outbound NATS data.
	triggerQueue := make(chan *busQueueEntry, 50)

	return &natsBus{apiConsumer, id, prefix, ec, triggerQueue}, nil
}

func (bus *natsBus) Start() {
	// Start go routine to process outbount NATS publishing.
	go chPublish(bus.ec, bus.triggerQueue)
	bus.Trigger(BusManagerStartup, bus.id, "", nil, nil)
}

func (bus *natsBus) Publish(subject string, v interface{}) error {
	return bus.ec.Publish(subject, v)
}

func (bus *natsBus) Request(subject string, v interface{}, vPtr interface{}, timeout time.Duration) error {
	return bus.ec.Request(subject, v, vPtr, timeout)
}

func (bus *natsBus) Trigger(name, from, payload string, data interface{}, pipeline *Pipeline) (err error) {
	trigger := &BusTrigger{
		Id:      bus.id,
		Name:    name,
		From:    from,
		Payload: payload,
		Data:    data,
	}
	if pipeline != nil {
		trigger.Pipeline = pipeline.GetID()
	}
	entry := &busQueueEntry{BusSubjectTrigger(bus.prefix, name), trigger}
	select {
	case bus.triggerQueue <- entry:
		// sent ok
	default:
		log.Println("Failed to queue NATS event - queue full?")
		err = errors.New("NATS trigger queue full")
	}

	return err
}

func (bus *natsBus) PrefixSubject(sub string) string {
	return fmt.Sprintf("%s.%s", bus.prefix, sub)
}

func (bus *natsBus) Subscribe(subject string, cb nats.Handler) (*nats.Subscription, error) {
	return bus.ec.Subscribe(subject, cb)
}

func (bus *natsBus) BindRecvChan(subject string, channel interface{}) (*nats.Subscription, error) {
	return bus.ec.BindRecvChan(subject, channel)
}

func (bus *natsBus) BindSendChan(subject string, channel interface{}) error {
	return bus.ec.BindSendChan(subject, channel)
}

func (bus *natsBus) CreateSink(id string) (sink Sink) {
	sink = newNatsSink(bus, id)
	return
}

type busQueueEntry struct {
	subject string
	data    interface{}
}

func chPublish(ec *natsconnection.EncodedConnection, channel chan (*busQueueEntry)) {
	for {
		entry := <-channel
		err := ec.Publish(entry.subject, entry.data)
		if err != nil {
			log.Println("Failed to publish to NATS", entry.subject, err)
		}
	}
}

type natsSink struct {
	sync.RWMutex
	id         string
	bm         BusManager
	closed     bool
	SubjectOut string
	SubjectIn  string
	sub        *nats.Subscription
	sendQueue  chan *DataSinkOutgoing
}

func newNatsSink(bm BusManager, id string) *natsSink {
	sink := &natsSink{
		id:         id,
		bm:         bm,
		SubjectOut: bm.PrefixSubject(fmt.Sprintf("sink.%s.out", id)),
		SubjectIn:  bm.PrefixSubject(fmt.Sprintf("sink.%s.in", id)),
	}

	sink.sendQueue = make(chan *DataSinkOutgoing, 100)
	bm.BindSendChan(sink.SubjectOut, sink.sendQueue)

	return sink
}

func (sink *natsSink) Write(outgoing *DataSinkOutgoing) (err error) {
	if sink.Enabled() {
		log.Println("Sending via NATS sink", sink.SubjectOut, outgoing)
		sink.sendQueue <- outgoing
	}
	return err
}

func (sink *natsSink) Enabled() bool {
	sink.RLock()
	defer sink.RUnlock()
	return sink.closed == false
}

func (sink *natsSink) Close() {
	sink.Lock()
	defer sink.Unlock()
	if sink.sub != nil {
		err := sink.sub.Unsubscribe()
		if err != nil {
			log.Println("Failed to unsubscribe NATS sink", err)
		} else {
			sink.sub = nil
		}
	}
	sink.closed = true
}

func (sink *natsSink) Export() *DataSink {
	return &DataSink{
		SubjectOut: sink.SubjectOut,
		SubjectIn:  sink.SubjectIn,
	}
}

func (sink *natsSink) BindRecvChan(channel interface{}) (*nats.Subscription, error) {
	sink.Lock()
	defer sink.Unlock()
	if sink.sub != nil {
		sink.sub.Unsubscribe()
		sink.sub = nil
	}
	sub, err := sink.bm.BindRecvChan(sink.SubjectIn, channel)
	if err != nil {
		return nil, err
	}
	sink.sub = sub
	return sub, nil
}
