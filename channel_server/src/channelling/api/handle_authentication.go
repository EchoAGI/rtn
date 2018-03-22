package api

import (
	"log"

	"channelling"
)

func (api *channellingAPI) HandleAuthentication(session *channelling.Session, st *channelling.SessionToken) (*channelling.DataSelf, error) {
	if err := api.SessionManager.Authenticate(session, st, ""); err != nil {
		log.Println("Authentication failed", err, st.Userid, st.Nonce)
		return nil, err
	}

	log.Println("Authentication success", session.Userid())
	self, err := api.HandleSelf(session)
	if err == nil {
		session.BroadcastStatus()
	}

	return self, err
}
