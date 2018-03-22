package logic


type HubManager interface {
	AddClient(client *Client) bool
	RemoteClient(client *Client)
	GetClient(cid string) *Client
}


type hubManager struct {
	clients map[string]*Client
}


func NewHubManager() *hubManager {
	return &hubManager{
	}
}


func (hub *hubManager) AddClient(client *Client) bool {

	if _, ok := hub.clients[client.cid]; ok {
		//udpate
	} else {
		hub.clients[client.cid] = client
		client.Start()
	}

	return true
}

func (hub *hubManager) RemoteClient(client *Client) {
	delete(hub.clients, client.cid)

	return
}

func (hub *hubManager) GetClient(cid string) *Client {
	return hub.clients[cid]
}