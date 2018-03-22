package channelling

type SessionCreateRequest struct {
	Id           string
	Session      *DataSession
	Room         *DataRoom
	SetAsDefault bool
}

type DataSink struct {
	SubjectOut string `json:subject_out"`
	SubjectIn  string `json:subject_in"`
}

type DataSinkOutgoing struct {
	Outgoing   *DataOutgoing
	ToUserid   string
	FromUserid string
	Pipe       string `json:",omitempty"`
}
