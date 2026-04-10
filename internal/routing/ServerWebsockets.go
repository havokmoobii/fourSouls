package routing

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type ServerConfig struct {
	Clients map[string]*websocket.Conn
	Rooms   []room
}

type room struct {
	clients map[string]*websocket.Conn
	state   string
}

func (cfg *ServerConfig) HandleConnect(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	roomNumber, err := strconv.Atoi(r.Header["Room"][0])
	if err != nil {
		http.Error(w, "Malformed header", http.StatusBadRequest)
		return
	}

	log.Println("Recieved a connection request from", username, "to join room", roomNumber)

	_, usernameTaken := cfg.Clients[username]
	if usernameTaken {
		log.Println("Responding to a failed connection request: Username is taken")
		http.Error(w, "Username is taken", http.StatusBadRequest)
		return
	}

	if len(cfg.Rooms[roomNumber-1].clients) > 3 {
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

	cfg.Rooms[roomNumber-1].clients[username] = conn

	log.Println("Client Successfully Connected")

	for {
		var pst Post
		err = conn.ReadJSON(&pst)
		if err != nil {
			log.Println("Read error:", err)
			log.Println("Removing disconnected user", username, "from room", roomNumber)
			delete(cfg.Rooms[roomNumber-1].clients, username)
			break
		}

		if pst.Kind == PostPlayerJoined {
			cfg.sendLobbyUpdate()
		}

		log.Println(pst)

		for _, client := range cfg.Rooms[roomNumber-1].clients {
			if err := client.WriteJSON(pst); err != nil {
				log.Println("Write error:", err)
				break
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (cfg *ServerConfig) sendLobbyUpdate() {
	log.Println("Sending a lobby update request to all clients")
	for _, room := range cfg.Rooms {
		for _, client := range room.clients {
			update := Post{
				Kind: PostLobbyUpdate,
				Msg: Message{
					Sender: "Server",
				},
			}

			if err := client.WriteJSON(update); err != nil {
				log.Println("Write error:", err)
				break
			}

		}
	}
}
