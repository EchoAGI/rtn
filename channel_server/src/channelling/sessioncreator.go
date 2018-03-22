package channelling

type SessionCreator interface {
	CreateSession(st *SessionToken, userid string) *Session
}
