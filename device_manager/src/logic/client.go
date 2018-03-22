package logic

import(
	"bus"
)

type Client struct {
	cid string
	clientAPI ClientAPI
	codec Codec
	topic string
	busManager bus.BusManager
}


func NewClient(cid string, clientAPI ClientAPI, codec Codec, topic string, busManager bus.BusManager) *Client {
	return &Client{
		cid : cid,
		clientAPI : clientAPI,
		codec : codec,
		topic : topic,
		busManager : busManager,
	}
}


func (c *Client) onMessage(topic string, in []byte) {
	var dataMsg DataMessage
	c.codec.Decode(in, &dataMsg)

	switch dataMsg.Type {
	case MSG_CLIENT_MSG:
		c.clientAPI.HandleMessage(c, &dataMsg.DataIncoming)
	}
}



func(c *Client) Start() {
	c.busManager.Subscribe(c.topic, c.onMessage)
}

func(c *Client) Stop() {
	c.busManager.UnSubscribe(c.topic)
}




func (c *Client) OnConnect(*Client, ) (interface{}, error) {
	var obj interface{}
	return obj, nil
}

func (c *Client) OnDisconnect(*Client) {

}

func (c *Client) OnIncoming(*Client, *DataIncoming) (interface{}, error) {
	var obj interface{}
	return obj, nil
}
