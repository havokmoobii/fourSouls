package routing

import (
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/websocket"
)

type ServerConfig struct {
	Clients     map[string]*websocket.Conn
	ChatClients map[string]*websocket.Conn
}

type ServerStatusResp struct {
	Games []Games
}

type Games struct {
	State string
	Users []string
}

func (cfg *ServerConfig) HandleStatus(w http.ResponseWriter, r *http.Request) {
	status := ServerStatusResp{}

	status.Games = append(status.Games, Games{})

	status.Games[0].State = "Waiting to Start"

	for username, _ := range cfg.Clients {
		status.Games[0].Users = append(status.Games[0].Users, username)
	}

	log.Println("Responding to a Status Request")

	respondWithJSON(w, http.StatusOK, status)
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

	if len(cfg.Clients) > 3 {
		http.Error(w, "Only 4 players can play per game", http.StatusBadRequest)
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

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

