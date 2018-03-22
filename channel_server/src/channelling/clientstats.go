package channelling

type ClientStats interface {
	ClientInfo(details bool) (int, map[string]*DataSession, map[string]string)
}
