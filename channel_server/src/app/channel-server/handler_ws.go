package main

import (
	"log"
	"net/http"

	"channelling"
	"channelling/server"

	"github.com/gorilla/websocket"
)

const (
	wsReadBufSize  = 1024
	wsWriteBufSize = 1024
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  wsReadBufSize,
		WriteBufferSize: wsWriteBufSize,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all connections by default to keep backwards
			// compatibility, but we should really check the Origin header instead!
			//
			// NOTE: We can omit "CheckOrigin" if the host in Origin must be the same as the host of the request (which is probably always the case).
			return true
		},
	}
)

func makeWSHandler(connectionCounter channelling.ConnectionCounter, sessionManager channelling.SessionManager, codec channelling.Codec, channellingAPI channelling.ChannellingAPI, users *server.Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate incoming request.
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Upgrade to Websocket mode.
		ws, err := upgrader.Upgrade(w, r, nil)
		if _, ok := err.(websocket.HandshakeError); ok {
			return
		} else if err != nil {
			log.Println(err)
			return
		}

		//如果开启账户模式, 要携带登录token
		r.ParseForm()
		token := r.FormValue("t")

		st := sessionManager.DecodeSessionToken(token)

		log.Printf("token: %s\n", token)
		log.Printf("st id: %s, sid: %s\n", st.Id, st.Sid)


		var userid string
		if users != nil {
			userid, _ = users.GetUserID(r)
			if userid == "" {
				userid = st.Userid
			}
		}

		// Create a new connection instance.
		session := sessionManager.CreateSession(st, userid)
		client := channelling.NewClient(codec, channellingAPI, session)
		conn := channelling.NewConnection(connectionCounter.CountConnection(), ws, client)

		// Start pumps (readPump blocks).
		go conn.WritePump()
		conn.ReadPump()
	}
}
