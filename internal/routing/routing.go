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
	ChatConn    *websocket.Conn
	CloseSignal bool
	Username    string
	GS          gamelogic.GameState
}

type Post struct {
	GS  gamelogic.GameState
	Msg Message
}

type Message struct {
	Sender    string
	Recipient string
	Body      string
}

func (cfg *ClientConfig) CheckServer() error {
	url := "http://localhost:1337/status"

	fmt.Println("\nChecking server for existing games...")

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

func (cfg *ClientConfig) Connect() error {
	for {
		username, err := gamelogic.ClientWelcome()

		url := fmt.Sprintf("ws://localhost:1337/connect/%s", username)
		chatUrl := fmt.Sprintf("ws://localhost:1337/chat/connect/%s", username)

		fmt.Println("\nConnecting to server...")

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

		ChatConn, dialResp, err := websocket.DefaultDialer.Dial(chatUrl, nil)
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
		cfg.ChatConn = ChatConn
		cfg.Username = username

		fmt.Println("Success!")

		// Have the chat Connection post that the player has joined the game here 

		return nil
	}
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
	for {
		var pst Post
		err := cfg.Conn.ReadJSON(&pst)
		if cfg.CloseSignal {
			return
		}
		if err != nil {
			fmt.Println("Read error:", err)
		}
		fmt.Println("Message Received")
		fmt.Println("> ")
		cfg.GS = pst.GS
	}
}

func (cfg *ClientConfig) ChatPost(msg string) error {
	err := cfg.ChatConn.WriteJSON(Message{
		Sender: cfg.Username,
		Body:   msg,
	})
	if err != nil {
		fmt.Println("Write error:", err)
		return err
	}

	return nil
}

// Can store everyone's usernames in GS to check if a recipient username is valid later.
func (cfg *ClientConfig) ChatDM(recipient, msg string) error {
	err := cfg.ChatConn.WriteJSON(Message{
		Sender: cfg.Username,
		Recipient: recipient,
		Body:   msg,
	})
	if err != nil {
		fmt.Println("Write error:", err)
		return err
	}

	return nil
}

func (cfg *ClientConfig) ReceiveChatPost() {
	for {
		var msg Message
		err := cfg.ChatConn.ReadJSON(&msg)
		if cfg.CloseSignal {
			return
		}
		if err != nil {
			fmt.Println("Read error:", err)
		}
		if msg.Recipient == "" {
			if cfg.Username == msg.Sender {
				fmt.Printf("\n<%s> %s\n\n> ", msg.Sender, msg.Body)
			} else {
				fmt.Printf("\n\n<%s> %s\n\n> ", msg.Sender, msg.Body)
			}
		} else {
			if cfg.Username == msg.Recipient {
				fmt.Printf("\n\n<<From:%s>> %s\n\n> ", msg.Sender, msg.Body)
			}
			if cfg.Username == msg.Sender {
				fmt.Printf("\n\n<<To:%s>> %s\n\n> ", msg.Recipient, msg.Body)
			}
		}
		
	}
}

// Rename this to ReceivePost when done
func (cfg *ClientConfig) ReceiveCombinedPost() {
// Maybe have a second loop before the game starts for different behavior?
// Look at message contents. If not "" then it is a message post. Maybe leave room to handle a post that is multiple types?
}
