package api

import (
	"channelling"
)

func (api *channellingAPI) HandleSessions(session *channelling.Session, sessions *channelling.DataSessionsRequest) (*channelling.DataSessions, error) {
	switch sessions.Type {
	case "contact":
		//guard 是否开启 联系人模块
		if !api.config.WithModule("contacts") {
			return nil, channelling.NewDataError("contacts_not_enabled", "incoming contacts session request with contacts disabled")
		}

		//获取联系人信息
		userID, err := api.ContactManager.GetContactID(session, sessions.Token)
		if err != nil {
			return nil, err
		}

		//得到指定联系人的会话信息
		return &channelling.DataSessions{
			Type:     "Sessions",
			Users:    api.SessionManager.GetUserSessions(session, userID),
			Sessions: sessions,
		}, nil
	case "session":
		id, err := session.DecodeAttestation(sessions.Token)
		if err != nil {
			return nil, channelling.NewDataError("bad_attestation", err.Error())
		}

		session, ok := api.Unicaster.GetSession(id)
		if !ok {
			return nil, channelling.NewDataError("no_such_session", "cannot retrieve session")
		}

		return &channelling.DataSessions{
			Type:     "Sessions",
			Users:    []*channelling.DataSession{session.Data()},
			Sessions: sessions,
		}, nil
	default:
		return nil, channelling.NewDataError("bad_request", "unknown sessions request type")
	}
}
