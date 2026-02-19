package routing

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

type ServerConfig struct {
	Clients     map[string]*websocket.Conn
	ChatClients map[string]*websocket.Conn
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func (cfg *ServerConfig) HandleConnections(w http.ResponseWriter, r *http.Request) {	
	username := r.PathValue("username")
	
	_, usernameTaken := cfg.Clients[username]
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
	
	cfg.Clients[username] = conn

	log.Println("Client Successfully Connected")

	for {
		var msg interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			delete(cfg.Clients, username)
			break
		}

		log.Println(msg)

		for _, client := range cfg.Clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Println("Write error:", err)
				break
			}
		}
	}
}

func (cfg *ServerConfig) HandleChatConnections(w http.ResponseWriter, r *http.Request) {	
	username := r.PathValue("username")
	
	_, usernameTaken := cfg.ChatClients[username]
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
	
	cfg.ChatClients[username] = conn

	log.Println("Chat Client Successfully Connected")

	for {
		var msg interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			delete(cfg.ChatClients, username)
			break
		}

		log.Println(msg)

		for _, client := range cfg.ChatClients {
			if err := client.WriteJSON(msg); err != nil {
				log.Println("Write error:", err)
				break
			}
		}
	}
}

