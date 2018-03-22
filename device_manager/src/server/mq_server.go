package server

import (
	log "github.com/sirupsen/logrus"
	"logic"
	"api"
	"bus"
)


type MQServer struct {
	topic string
	api  api.API
	busManager bus.BusManager
	hubManager logic.HubManager
	codec logic.Codec
}

func NewMQServer(api api.API, busManager bus.BusManager, hubManager logic.HubManager, codec logic.Codec) *MQServer {
	return &MQServer{
		api : api,
		busManager : busManager,
		hubManager : hubManager,
		codec : codec,
	}
}


func(s *MQServer) serveSubscribe() {
	s.busManager.Subscribe(s.topic, func(topic string, payload []byte) {
		s.process(topic, payload)
	})
}

func (s *MQServer) process(topic string, payload []byte) {
	var dataMsg logic.DataMessage
	s.codec.Decode(payload, &dataMsg)

	switch dataMsg.Type {
	case logic.MSG_REGISTER_DEVICE:
		cid := dataMsg.RegisterDevice.DeviceId

		clientApi := logic.NewClientAPI()
		client := logic.NewClient(cid, clientApi, s.codec, dataMsg.RegisterDevice.Topic, s.busManager)
		b := s.hubManager.AddClient(client)
		if b {
			s.api.HandleRegisterDevice(client, &dataMsg.RegisterDevice)
		}
	case logic.MSG_UNREGISTER_DEVICE:
		cid := dataMsg.RegisterDevice.DeviceId

		s.hubManager.GetClient(cid)

		client := s.hubManager.GetClient(cid)
		if client == nil {
			return
		}

		s.api.HandleUnRegisterDevice(client, &dataMsg.UnRegisterDevice)
	default:
		log.Println("OnText unhandled message type", dataMsg.Type)
	}
}

func (s *MQServer) Serve() {
	s.serveSubscribe()
}