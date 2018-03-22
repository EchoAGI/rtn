package channelling

const (
	RoomTypeConference = "Conference"
	RoomTypeRoom       = "Room"
)

type ChannellingAPI interface {
	OnConnect(*Client, *Session) (interface{}, error)
	OnDisconnect(*Client, *Session)
	OnIncoming(Sender, *Session, *DataIncoming) (interface{}, error)
	OnIncomingProcessed(Sender, *Session, *DataIncoming, interface{}, error)
}

type ChannellingAPIConsumer interface {
	SetChannellingAPI(ChannellingAPI)
	GetChannellingAPI() ChannellingAPI
}

type channellingAPIConsumer struct {
	ChannellingAPI ChannellingAPI
}

func NewChannellingAPIConsumer() ChannellingAPIConsumer {
	return &channellingAPIConsumer{}
}

func (c *channellingAPIConsumer) SetChannellingAPI(api ChannellingAPI) {
	c.ChannellingAPI = api
}

func (c *channellingAPIConsumer) GetChannellingAPI() ChannellingAPI {
	return c.ChannellingAPI
}
