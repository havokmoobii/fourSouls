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

func (cfg *serverConfig) handleConnections(w http.ResponseWriter, r *http.Request) {	
	username := r.PathValue("username")
	
	_, usernameTaken := cfg.clients[username]
	if usernameTaken {
		http.Error(w, "Username is taken", http.StatusBadRequest)
		return
	}
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()
	
	cfg.clients[username] = conn

	log.Println("Client Successfully Connected")

	for {
		var msg interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Println(msg)

		for _, client := range cfg.clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Println("Write error:", err)
				break
			}
		}
	}
}

func (cfg *serverConfig) handleChatConnections(w http.ResponseWriter, r *http.Request) {	
	username := r.PathValue("username")
	
	_, usernameTaken := cfg.chatClients[username]
	if usernameTaken {
		http.Error(w, "Username is taken", http.StatusBadRequest)
		return
	}
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()
	
	cfg.chatClients[username] = conn

	log.Println("Chat Client Successfully Connected")

	for {
		var msg interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Println(msg)

		for _, client := range cfg.chatClients {
			if err := client.WriteJSON(msg); err != nil {
				log.Println("Write error:", err)
				break
			}
		}
	}
}

