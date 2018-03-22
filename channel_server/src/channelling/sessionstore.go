package channelling

type SessionStore interface {
	GetSession(id string) (session *Session, ok bool)
}
