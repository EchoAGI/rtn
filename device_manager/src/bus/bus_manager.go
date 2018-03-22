package bus





type MQFunc func(topic string, message []byte)



type BusManager interface {
	Subscribe(topic string, fun MQFunc)
	UnSubscribe(topic string)
	Publish(topic string, payload []byte) error
}







type natsBusManager struct {

}