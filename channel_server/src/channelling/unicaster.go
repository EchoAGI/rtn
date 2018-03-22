package channelling

type Unicaster interface {
	SessionStore
	OnConnect(*Client, *Session)
	OnDisconnect(*Client, *Session)
	Unicast(to string, outgoing *DataOutgoing, pipeline *Pipeline)
}
