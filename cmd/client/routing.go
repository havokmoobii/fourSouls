package main

import (
	"fmt"
	"io"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/havokmoobii/fourSouls/internal/gamelogic"
)

type clientConfig struct {
	conn     *websocket.Conn
	chatConn *websocket.Conn
	username string
	gs       gamelogic.GameState
}

type message struct {
	Sender    string
	Recipient string
	Body      string
}

func (cfg *clientConfig) connect() error {
	for {
		username, err := gamelogic.ClientWelcome()

		url := fmt.Sprintf("ws://localhost:1337/connect/%s", username)
		chatUrl := fmt.Sprintf("ws://localhost:1337/chat/connect/%s", username)

		fmt.Println("Connecting to server...")

		conn, dialResp, err := websocket.DefaultDialer.Dial(url, nil)
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

		chatConn, dialResp, err := websocket.DefaultDialer.Dial(chatUrl, nil)
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

		cfg.conn = conn
		cfg.chatConn = chatConn
		cfg.username = username

		// Have the chat connection post that the player has joined the game here 

		return nil
	}
}

func (cfg *clientConfig) post(msg gamelogic.GameState) error {
	err := cfg.conn.WriteJSON(msg)
	if err != nil {
		fmt.Println("Write error:", err)
		return err
	}

	return nil
}

func (cfg *clientConfig) receivePost() {
	for {
		var data gamelogic.GameState
		err := cfg.conn.ReadJSON(&data)
		if err != nil {
			fmt.Println("Read error:", err)
		}
		fmt.Println("Message Received")
		cfg.gs = data
	}
}

func (cfg *clientConfig) chatPost(msg string) error {
	err := cfg.chatConn.WriteJSON(message{
		Sender: cfg.username,
		Body:   msg,
	})
	if err != nil {
		fmt.Println("Write error:", err)
		return err
	}

	return nil
}

// Can store everyone's usernames in gs to check if a recipient username is valid later.
func (cfg *clientConfig) chatDM(recipient, msg string) error {
	err := cfg.chatConn.WriteJSON(message{
		Sender: cfg.username,
		Recipient: recipient,
		Body:   msg,
	})
	if err != nil {
		fmt.Println("Write error:", err)
		return err
	}

	return nil
}

func (cfg *clientConfig) receiveChatPost() {
	for {
		var msg message
		err := cfg.chatConn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Read error:", err)
		}
		if msg.Recipient == "" {
			fmt.Printf("\n<%s> %s\n>", msg.Sender, msg.Body)
		} else {
			if cfg.username == msg.Recipient {
				fmt.Printf("\n<<From:%s>> %s\n>", msg.Sender, msg.Body)
			}
			if cfg.username == msg.Sender {
				fmt.Printf("\n<<To:%s>> %s\n>", msg.Recipient, msg.Body)
			}
		}
		
	}
}
