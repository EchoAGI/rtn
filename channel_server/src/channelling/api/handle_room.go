package api

import (
	"channelling"
)

func (api *channellingAPI) HandleRoom(session *channelling.Session, room *channelling.DataRoom) (*channelling.DataRoom, error) {
	room, err := api.RoomStatusManager.UpdateRoom(session, room)
	if err == nil {
		session.Broadcast(room)
	}

	return room, err
}

func (api *channellingAPI) RoomProcessed(sender channelling.Sender, session *channelling.Session, msg *channelling.DataIncoming, reply interface{}, err error) {
	if err == nil {
		api.SendConferenceRoomUpdate(session)
	}
}

func (api *channellingAPI) SendConferenceRoomUpdate(session *channelling.Session) {
	// If user joined a server-managed conference room, send list of session ids to all participants.
	if room, ok := api.RoomStatusManager.Get(session.Roomid); ok && room.GetType() == channelling.RoomTypeConference {
		if sessionids := room.SessionIDs(); len(sessionids) > 1 {
			cid := session.Roomid
			session.Broadcaster.Broadcast("", session.Roomid, &channelling.DataOutgoing{
				To: cid,
				Data: &channelling.DataConference{
					Type:       "Conference",
					Id:         cid,
					Conference: sessionids,
				},
			})
		}
	}
}
