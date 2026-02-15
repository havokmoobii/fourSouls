package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func handlerWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
	}
	defer conn.Close()

	log.Println("Client Successfully Connected")

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Println(string(msg))

		if err := conn.WriteMessage(messageType, msg); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}