package server

import (
	"log"
	"net/http"
	"strings"

	"channelling"
)

type Token struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
}

type Tokens struct {
	Provider channelling.TokenProvider
}

func (tokens Tokens) Post(request *http.Request) (int, interface{}, http.Header) {

	auth := request.Form.Get("a")

	if len(auth) > 100 {
		return 413, NewApiError("auth_too_large", "Auth too large"), http.Header{"Content-Type": {"application/json"}}
	}

	valid := tokens.Provider(strings.ToLower(auth))

	if valid != "" {
		log.Printf("Good incoming token request: %s\n", auth)
		return 200, &Token{Token: valid, Success: true}, http.Header{"Content-Type": {"application/json"}}
	}
	log.Printf("Wrong incoming token request: %s\n", auth)
	return 403, NewApiError("invalid_token", "Invalid token"), http.Header{"Content-Type": {"application/json"}}

}
