package channelling

type UserStore interface {
	GetUser(id string) (user *User, ok bool)
}
