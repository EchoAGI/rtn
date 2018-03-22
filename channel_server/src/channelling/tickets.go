package channelling

import (
	"crypto/aes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"randomstring"

	"github.com/gorilla/securecookie"
	log "github.com/sirupsen/logrus"
)

var (
	// Can be set from tests to disable some log outputs.
	silentOutput = false
)

type SessionValidator interface {
	Realm() string
	ValidateSession(string, string) bool
}

type SessionEncoder interface {
	EncodeSessionToken(*Session) (string, error)
	EncodeSessionUserID(*Session) string
}

type Tickets interface {
	SessionValidator
	SessionEncoder
	DecodeSessionToken(token string) (st *SessionToken)
	FakeSessionToken(userid string) *SessionToken
}

type tickets struct {
	*securecookie.SecureCookie
	realm            string
	tokenName        string
	encryptionSecret []byte
}

func NewTickets(sessionSecret, encryptionSecret []byte, realm string) Tickets {
	tickets := &tickets{
		nil,
		realm,
		fmt.Sprintf("token@%s", realm),
		encryptionSecret,
	}
	tickets.SecureCookie = securecookie.New(sessionSecret, encryptionSecret)
	tickets.MaxAge(86400 * 30) // 30 days
	tickets.HashFunc(sha256.New)
	tickets.BlockFunc(aes.NewCipher)

	return tickets
}

func (tickets *tickets) Realm() string {
	return tickets.realm
}

func reverseBase64String(s string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	for i, j := 0, len(decoded)-1; i < j; i, j = i+1, j-1 {
		decoded[i], decoded[j] = decoded[j], decoded[i]
	}
	return base64.URLEncoding.EncodeToString(decoded), nil
}

func (tickets *tickets) reverseSessionId(id string) string {
	reversedId, err := reverseBase64String(id)
	if err != nil {
		// This should never happen
		panic("Could not reverse " + id)
	}
	return reversedId
}

/**
 * 否则随机生成
 */
func (tickets *tickets) DecodeSessionToken(token string) (st *SessionToken) {
	var err error
	if token != "" {
		st = &SessionToken{}
		err = tickets.Decode(tickets.tokenName, token, st)

		log.WithFields(log.Fields{
			"tokenName" : tickets.tokenName,
			"userId" : st.Userid,
			"appid" : st.Appid,
			"appSecret" : st.AppSecret,
		}).Info("DecodeSessionToken")

		if err != nil {
			log.Error("Error while decoding session token", err)
		}
	}

	if st == nil || err != nil {
		sid := randomstring.NewRandomString(32)
		id, _ := tickets.Encode("id", sid)
		id = tickets.reverseSessionId(id)
		st = &SessionToken{Id: id, Sid: sid}
		if !silentOutput {
			log.Println("Created new session id", id)
		}
	}
	return
}

func (tickets *tickets) FakeSessionToken(userid string) (st *SessionToken) {
	sid := fmt.Sprintf("fake-%s", randomstring.NewRandomString(27))
	id, _ := tickets.Encode("id", sid)
	id = tickets.reverseSessionId(id)
	st = &SessionToken{Id: id, Sid: sid, Userid: userid}
	log.Println("Created new fake session id", st.Id)
	return
}

func (tickets *tickets) ValidateSession(id, sid string) bool {
	var decoded string
	reversedId, err := reverseBase64String(id)
	if err != nil {
		if !silentOutput {
			log.Println("Session format error", err, id, sid)
		}
		return false
	}

	if err := tickets.Decode("id", reversedId, &decoded); err != nil {
		if !silentOutput {
			log.Println("Session validation error", err, reversedId, sid)
		}
		return false
	}
	
	if decoded != sid {
		if !silentOutput {
			log.Println("Session validation failed", reversedId, sid)
		}
		return false
	}
	return true
}

func (tickets *tickets) EncodeSessionToken(session *Session) (string, error) {
	return tickets.Encode(tickets.tokenName, session.Token())
}

func (tickets *tickets) EncodeSessionUserID(session *Session) (suserid string) {
	if userid := session.Userid(); userid != "" {
		m := hmac.New(sha256.New, tickets.encryptionSecret)
		m.Write([]byte(userid))
		suserid = base64.StdEncoding.EncodeToString(m.Sum(nil))
	}
	return
}
