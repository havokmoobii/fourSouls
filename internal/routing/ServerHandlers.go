// Server-Side Connections

package routing

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

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

	for roomNumber, room := range cfg.Rooms {
		status.Games = append(status.Games, Game{})
		status.Games[roomNumber].State = room.state
		for username := range room.clients {
			status.Games[roomNumber].Users = append(status.Games[roomNumber].Users, username)
		}
	}

	respondWithJSON(w, http.StatusOK, status)
}

func (cfg *ServerConfig) HandleRoomsCreate(w http.ResponseWriter, r *http.Request) {
	roomToOccupy := -1

	for roomNumber, room := range cfg.Rooms {
		if len(room.clients) == 0 {
			roomToOccupy = roomNumber
			break
		}
	}

	if roomToOccupy == -1 {
		cfg.Rooms = append(cfg.Rooms, room{
			clients: make(map[string]*websocket.Conn),
			state:   "Waiting to Start",
		})
	}

	resp := RoomsPostResponse{
		RoomNumber: len(cfg.Rooms),
	}

	if roomToOccupy != -1 {
		resp.RoomNumber = roomToOccupy + 1
	}

	log.Println("Responding to a room creation request")

	respondWithJSON(w, http.StatusOK, resp)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
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
