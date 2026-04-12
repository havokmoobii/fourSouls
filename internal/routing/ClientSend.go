// Client-Side Connections

package routing

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/havokmoobii/fourSouls/internal/gamelogic"
)

type ClientConfig struct {
	client      http.Client
	Conn        *websocket.Conn
	StartSignal bool
	CloseSignal bool
	Username    string
	RoomNumber  int
	GS          gamelogic.GameState
}

func (cfg *ClientConfig) Connect(roomNumber string) error {
	roomNumberInt, err := strconv.Atoi(roomNumber)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("ws://localhost:1337/connect/%s", cfg.Username)

	headers := http.Header{}
	headers.Add("Room", roomNumber)

	fmt.Println("\nConnecting to server...")
	fmt.Println()

	Conn, dialResp, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		if dialResp != nil {
			defer dialResp.Body.Close()
			body, _ := io.ReadAll(dialResp.Body)

			fmt.Printf("\nHTTP Status: %d %s\n", dialResp.StatusCode, http.StatusText(dialResp.StatusCode))
			fmt.Printf("Server message: %s\n", string(body))

			return err
		}
		fmt.Println("Dial error:", err)
		return err
	}

	cfg.Conn = Conn
	cfg.RoomNumber = roomNumberInt

	fmt.Println("Success!")

	go cfg.ReceivePost()

	return nil
}

func (cfg *ClientConfig) CheckServer() error {
	url := "http://localhost:1337/rooms"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := cfg.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	status := ServerStatusResp{}
	err = json.Unmarshal(dat, &status)
	if err != nil {
		if resp.StatusCode == http.StatusNoContent {
			fmt.Println("\nThe lobby is empty.")
			return nil
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Printf("\nHTTP Status: %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
		fmt.Printf("Server message: %s\n", string(body))
		return err
	}

	for roomNumber, game := range status.Games {
		if len(game.Users) == 0 {
			continue
		}
		fmt.Printf("\nGameroom %v:\n", roomNumber+1)
		for _, user := range game.Users {
			fmt.Println(user)
		}
	}

	return nil
}

func (cfg *ClientConfig) CreateRoom() error {
	url := "http://localhost:1337/rooms"

	fmt.Println("\nCreating and joining new Gameroom...")

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	resp, err := cfg.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	roomsResp := RoomsPostResponse{}
	err = json.Unmarshal(dat, &roomsResp)
	if err != nil {
		return err
	}

	err = cfg.JoinRoom(strconv.Itoa(roomsResp.RoomNumber))
	if err != nil {
		fmt.Println("Unable to join newly created room:", err)
	}

	return nil
}

func (cfg *ClientConfig) JoinRoom(roomNumber string) error {
	err := cfg.Connect(roomNumber)
	if err != nil {
		fmt.Printf("Unable to join room %s: %v\n", roomNumber, err)
		return err
	}

	err = cfg.SendPost(Post{
		Kind: PostPlayerJoined,
		Msg: Message{
			Sender: cfg.Username,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (cfg *ClientConfig) SendPost(pst Post) error {
	err := cfg.Conn.WriteJSON(pst)
	if err != nil {
		fmt.Println("Write error:", err)
		return err
	}

	return nil
}

func (cfg *ClientConfig) printChat(msg Message) {
	// Normal chat message
	if msg.Recipient == "" {
		if cfg.Username == msg.Sender {
			fmt.Printf("\n<%s> %s\n\n> ", msg.Sender, msg.Body)
		} else {
			fmt.Printf("\n\n<%s> %s\n\n> ", msg.Sender, msg.Body)
		}
		// Direct message to another user
	} else {
		if cfg.Username == msg.Sender {
			fmt.Printf("\n<<To:%s>> %s\n\n> ", msg.Recipient, msg.Body)
		}
		if cfg.Username == msg.Recipient {
			fmt.Printf("\n\n<<From:%s>> %s\n\n> ", msg.Sender, msg.Body)
		}
	}
}
