package channelling

import (
	"log"

	"buffercache"
)


const (
	// TextMessage denotes a text data message. The text message payload is interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2
)

type Message struct {
	buffercache.Buffer
	FrameType int
}


type Sender interface {
	Index() uint64
	Send(*Message)
}

type Client struct {
	Connection
	Codec
	ChannellingAPI ChannellingAPI
	session        *Session
}

func NewClient(codec Codec, api ChannellingAPI, session *Session) *Client {
	return &Client{
		Codec:          codec,
		ChannellingAPI: api,
		session:        session,
	}
}

func (client *Client) OnConnect(conn Connection) {
	client.Connection = conn
	if reply, err := client.ChannellingAPI.OnConnect(client, client.session); err == nil {
		client.reply("", reply)
	} else {
		log.Println("OnConnect error", err)
	}
}

func (client *Client) OnDisconnect() {
	client.session.Close()
	client.ChannellingAPI.OnDisconnect(client, client.session)
}

func (client *Client) OnText(b buffercache.Buffer) {
	incoming, err := client.Codec.DecodeIncoming(b)
	if err != nil {
		log.Println("OnText error while processing incoming message", err)
		return
	}

	var reply interface{}
	if reply, err = client.ChannellingAPI.OnIncoming(client, client.session, incoming); err != nil {
		client.reply(incoming.Iid, err)
	} else if reply != nil {
		client.reply(incoming.Iid, reply)
	}
	client.ChannellingAPI.OnIncomingProcessed(client, client.session, incoming, reply, err)
}

func (client *Client) reply(iid string, m interface{}) {
	outgoing := &DataOutgoing{From: client.session.Id, Iid: iid, Data: m}
	if b, err := client.Codec.EncodeOutgoing(outgoing); err == nil {
		msg := &Message{b, TextMessage}
		client.Connection.Send(msg)
		b.Decref()
	}
}

func (client *Client) Session() *Session {
	return client.session
}

func (client *Client) ReplaceAndClose(oldClient *Client) {
	oldSession := oldClient.Session()
	client.session.Replace(oldSession)
	go func() {
		// 在另一个 routine 中关闭老的 session & client, 避免因为老 session 挂起等问题 阻塞新的 client.
		log.Printf("Closing obsolete client %d (replaced with %d) with id %s\n", oldClient.Index(), client.Index(), oldSession.Id)
		oldSession.Close()
		oldClient.Close()
	}()
}
