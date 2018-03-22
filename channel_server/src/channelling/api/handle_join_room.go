package api

import (
	"channelling"
)

/**
 * 进入房间
 */
func (api *channellingAPI) HandleJoinRoom(session *channelling.Session, dataJoinRoom *channelling.DataJoinRoom, sender channelling.Sender) (*channelling.DataWelcome, error) {
	// TODO(longsleep): Filter room id and user agent.
	session.Update(&channelling.SessionUpdate{Types: []string{"Ua"}, Ua: dataJoinRoom.Ua})

	// Compatibily for old clients.
	roomName := dataJoinRoom.Name
	if roomName == "" {
		roomName = dataJoinRoom.Id
	}

	room, err := session.JoinRoom(roomName, dataJoinRoom.Type, dataJoinRoom.Credentials, sender)
	if err != nil {
		return nil, err
	}

	return &channelling.DataWelcome{
		Type:  "Welcome",
		Room:  room,
		Users: api.RoomStatusManager.RoomUsers(session),
	}, nil
}

func (api *channellingAPI) JoinRoomProcessed(sender channelling.Sender, session *channelling.Session, msg *channelling.DataIncoming, reply interface{}, err error) {
	if err == nil {
		api.SendConferenceRoomUpdate(session)
	}
}