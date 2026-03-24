package routing

import (
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/websocket"
)

type ServerConfig struct {
	Clients     map[string]*websocket.Conn
	Rooms       []room
}

type room struct {
	clients     map[string]*websocket.Conn
	game        Game
}

type RoomsPostResponse struct {
	RoomNumber int
}

type ServerStatusResp struct {
	Games []Game
}

type Game struct {
	State string
	Users []string
}

func (cfg *ServerConfig) HandleRooms(w http.ResponseWriter, r *http.Request) {
	status := ServerStatusResp{}

	if len(cfg.Rooms) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	for _, room := range cfg.Rooms {
		status.Games = append(status.Games, room.game)
 	}

	log.Println("Responding to a Status Request")

	respondWithJSON(w, http.StatusOK, status)
}

func (cfg *ServerConfig) HandleRoomsCreate(w http.ResponseWriter, r *http.Request) {
	cfg.Rooms = append(cfg.Rooms, room{
		game: Game{
			State: "Waiting to Start",
		},
	})

	resp := RoomsPostResponse{
		RoomNumber: len(cfg.Rooms),
	}

	log.Println("Responding to a room creation request")

	respondWithJSON(w, http.StatusOK, resp)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func (cfg *ServerConfig) HandleConnect(w http.ResponseWriter, r *http.Request) {	
	username := r.PathValue("username")

	roomNumber := r.Header["Room"]

	log.Println("Recieved a connection request from", username, "to join room", roomNumber)
	
	_, usernameTaken := cfg.Clients[username]
	if usernameTaken {
		log.Println("Responding to a failed connection request: Username is taken")
		http.Error(w, "Username is taken", http.StatusBadRequest)
		return
	}

	if len(cfg.Clients) > 3 {
		log.Println("Responding to a failed connection request: Only 4 players can play per game")
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

