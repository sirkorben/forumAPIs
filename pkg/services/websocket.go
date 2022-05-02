package services

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func webSocket(w http.ResponseWriter, r *http.Request) {
	// 	log.Println("socket request")
	// 	/*
	// 		defer is used to ensure that a function call is performed later in a program’s execution,
	// 		usually for purposes of cleanup.
	// 	*/
	// 	defer func() {
	// 		err := recover()
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 		r.Body.Close()
	// 	}()
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("client connected.")
	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}
