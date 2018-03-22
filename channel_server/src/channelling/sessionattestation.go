package channelling

import (
	"time"
)

type SessionAttestation struct {
	refresh int64
	token   string
	s       *Session
}

func (sa *SessionAttestation) Update() (string, error) {
	token, err := sa.Encode()
	if err == nil {
		sa.token = token
		sa.refresh = time.Now().Unix() + 180 // expires after 3 minutes
	}
	return token, err
}

func (sa *SessionAttestation) Token() (token string) {
	if sa.refresh < time.Now().Unix() {
		token, _ = sa.Update()
	} else {
		token = sa.token
	}
	return
}

func (sa *SessionAttestation) Encode() (string, error) {
	return sa.s.attestations.Encode("attestation", sa.s.Id)
}

func (sa *SessionAttestation) Decode(token string) (string, error) {
	var id string
	err := sa.s.attestations.Decode("attestation", token, &id)
	return id, err
}
