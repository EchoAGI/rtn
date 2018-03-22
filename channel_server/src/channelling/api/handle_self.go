package api

import (
	"log"

	"channelling"
)

func (api *channellingAPI) HandleSelf(session *channelling.Session) (*channelling.DataSelf, error) {
	token, err := api.SessionEncoder.EncodeSessionToken(session)
	if err != nil {
		log.Println("Error in OnRegister", err)
		return nil, err
	}

	log.Println("Created new session token", len(token), token)
	self := &channelling.DataSelf{
		Type:       "Self",
		Id:         session.Id,
		Sid:        session.Sid,
		Userid:     session.Userid(),
		Suserid:    api.SessionEncoder.EncodeSessionUserID(session),
		Token:      token,
		Version:    api.config.Version,
		ApiVersion: apiVersion,
		Turn:       api.TurnDataCreator.CreateTurnData(session),
		Stun:       api.config.StunURIs,
	}
	api.BusManager.Trigger(channelling.BusManagerSession, session.Id, session.Userid(), nil, nil)

	return self, nil
}
