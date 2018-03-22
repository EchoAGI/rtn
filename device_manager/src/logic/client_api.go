package logic


type ClientAPI interface {
	/*
	OnConnect(*Client, ) (interface{}, error)
	OnDisconnect(*Client)
	OnIncoming(*Client, *DataIncoming) (interface{}, error)
	*/
	HandleMessage(*Client, *DataIncoming)
}


type ClientAPIImpl struct {
}

func NewClientAPI() *ClientAPIImpl {
	return &ClientAPIImpl {

	}
}

func(c *ClientAPIImpl) HandleMessage(*Client, *DataIncoming) {

}