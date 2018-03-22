package channelling

type SessionUpdate struct {
	Types  []string
	Ua     string			//user agent
	Prio   int
	Status interface{}
}
