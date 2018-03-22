package channelling

import (
	"log"
	"sort"
	"sync"
)

type User struct {
	Id           string
	sessionTable map[string]*Session
	mutex        sync.RWMutex
}

func NewUser(id string) *User {
	user := &User{
		Id:           id,
		sessionTable: make(map[string]*Session),
	}

	return user
}

// AddSession adds a session to the session table and returns true if
// s is the first session.
func (u *User) AddSession(s *Session) bool {
	first := false
	u.mutex.Lock()
	u.sessionTable[s.Id] = s
	if len(u.sessionTable) == 1 {
		log.Println("First session registered for user", u.Id)
		first = true
	}
	u.mutex.Unlock()

	return first
}

// RemoveSession removes a session from the session table abd returns
// true if no session is left left.
func (u *User) RemoveSession(sessionID string) bool {
	last := false
	u.mutex.Lock()
	delete(u.sessionTable, sessionID)
	if len(u.sessionTable) == 0 {
		log.Println("Last session unregistered for user", u.Id)
		last = true
	}
	u.mutex.Unlock()

	return last
}

func (u *User) Data() *DataUser {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	return &DataUser{
		Id:       u.Id,
		Sessions: len(u.sessionTable),
	}
}

func (u *User) SubscribeSessions(from *Session) []*DataSession {
	sessions := make([]*DataSession, 0, len(u.sessionTable))
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	for _, session := range u.sessionTable {
		// TODO(longsleep): This does lots of locks - check if these can be streamlined.
		from.Subscribe(session)
		sessions = append(sessions, session.Data())
	}
	sort.Sort(ByPrioAndStamp(sessions))

	return sessions
}

type ByPrioAndStamp []*DataSession

func (a ByPrioAndStamp) Len() int {
	return len(a)
}

func (a ByPrioAndStamp) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByPrioAndStamp) Less(i, j int) bool {
	if a[i].Prio < a[j].Prio {
		return true
	}
	if a[i].Prio == a[j].Prio {
		return a[i].stamp < a[j].stamp
	}

	return false
}
