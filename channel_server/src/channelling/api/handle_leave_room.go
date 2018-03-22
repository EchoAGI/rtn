package api

import (
	"channelling"
)

func (api *channellingAPI) HandleLeave(session *channelling.Session) error {
	session.LeaveRoom()

	return nil
}