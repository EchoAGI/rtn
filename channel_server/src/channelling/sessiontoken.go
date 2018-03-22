package channelling

type SessionToken struct {
	Id     		string 	// Public session id.
	Sid    		string 	// Secret session id.
	Userid 		string 	// Public user id.
	Appid		string	// Appid
	AppSecret	string	// AppSecret
	Nonce  		string 	`json:"Nonce,omitempty"` // User autentication nonce.
}