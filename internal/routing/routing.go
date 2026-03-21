package routing

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/havokmoobii/fourSouls/internal/gamelogic"
)

type ClientConfig struct {
	client      http.Client
	Conn        *websocket.Conn
	StartSignal bool
	CloseSignal bool
	Username    string
	GameRoomNum int
	GS          gamelogic.GameState
}

type PostKind int

const (
	PostPlayerJoined PostKind = iota
	PostGameStart 
	PostStateUpdate
	PostChat
)

type Post struct {
	Kind        PostKind
	GameRoomNum int
	GS          gamelogic.GameState
	Msg         Message
}

type Message struct {
	Sender    string
	Recipient string
	Body      string
}

func (cfg *ClientConfig) Connect() error {
	for {
		url := fmt.Sprintf("ws://localhost:1337/connect/%s", cfg.Username)

		fmt.Println("\nConnecting to server...\n")

		Conn, dialResp, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			if dialResp != nil {
				body, _ := io.ReadAll(dialResp.Body)
				dialResp.Body.Close()

				fmt.Printf("\nHTTP Status: %d %s\n", dialResp.StatusCode, http.StatusText(dialResp.StatusCode))
				fmt.Printf("Server message: %s\n", string(body))
				continue
			}
			fmt.Println("Dial error:", err)
			return err
		}

		cfg.Conn = Conn

		fmt.Println("Success!")

		go cfg.ReceivePost()

		return nil
	}
}

func (cfg *ClientConfig) CheckServer() error {
	url := "http://localhost:1337/status"

	if cfg.GameRoomNum == 0 {
		fmt.Println("\nChecking server for existing games...")
	}

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
		return err
	}

	if len(status.Games[0].Users) == 0 {
		fmt.Println("\nThe lobby is empty.\n")
		return nil
	}

	fmt.Println("\nGameroom 1:", status.Games[0].State)
	for _, user := range status.Games[0].Users {
		fmt.Println(user)
	}
	fmt.Println()
	
	return nil
}

func (cfg *ClientConfig) CreateRoom() error {
	url := "http://localhost:1337/room"

	fmt.Println("\nCreating and joining new Gameroom...")

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
		return err
	}

	if len(status.Games[0].Users) == 0 {
		fmt.Println("\nThe lobby is empty.\n")
		return nil
	}

	fmt.Println("\nGameroom 1:", status.Games[0].State)
	for _, user := range status.Games[0].Users {
		fmt.Println(user)
	}
	fmt.Println()
	
	return nil
}

func (cfg *ClientConfig) JoinRoom(roomNumber int) error {





	// Next time: Figure out parameters to get the server to see room numbers and make seperate games.
	// Failing that we can do with the current implementation of all clients recieving all messages and just ingnoring the ones from other games.
	// that seems kinda sloppy though...




	err := cfg.Connect()
			if err != nil {
				fmt.Println("Unable to connect to server:", err)
				return err
			}

	cfg.GameRoomNum = roomNumber
	
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

func (cfg *ClientConfig) ReceivePost() {
	// Maybe have a second loop before the game starts for different behavior?
	for {
		var pst Post
		err := cfg.Conn.ReadJSON(&pst)

		// Only Recieve messages from your game.
		if pst.GameRoomNum != cfg.GameRoomNum{
			continue
		}

		if cfg.CloseSignal {
			return
		}
		if err != nil {
			fmt.Println("Read error:", err)
		}
		if pst.Kind == PostPlayerJoined {
			fmt.Println("Someone has joined the lobby!")
		}
		if pst.Kind == PostGameStart {
			cfg.StartSignal = true
		}
		if pst.Kind == PostChat {
			cfg.printChat(pst.Msg)
		}
		if pst.Kind == PostStateUpdate {
			if cfg.Username == pst.Msg.Sender{
				fmt.Print("\nGameState updated!\n\n> ")
			} else {
				fmt.Print("\n\nGameState updated!\n\n> ")
			}
			cfg.GS = pst.GS
		}	
	}
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
