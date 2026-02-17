package main

// Next time: Set up a repl loop and give the clients options.
// I want to set up chat, but not sure how with the channel blocking implementation

import (
	"fmt"
	"io"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/havokmoobii/fourSouls/internal/gamelogic"
)

func connect() (*websocket.Conn, error) {
	for {
		username, err := gamelogic.ClientWelcome()

		url := fmt.Sprintf("ws://localhost:1337/connect/%s", username)

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
			return nil, err
		}
		return conn, nil
	}
}

func post(conn *websocket.Conn, msg interface{}) error {
	err := conn.WriteJSON(msg)
	if err != nil {
		fmt.Println("Write error:", err)
		return err
	}

	return nil
}

func main() {
	conn, err := connect()
	if err != nil {
		// Error is handled in the function, so we can simply return.
		return
	}
	defer conn.Close()

	fmt.Println("Success!")

	resp := make(chan interface{})
	go func() {
		defer close(resp)
		for {
			var payload interface{}
			err = conn.ReadJSON(&payload)
			resp <- payload
		}
	}()

	gs := gamelogic.GameState{
		Player: "HavokMoobii",
	}

	err = post(conn, gs)
	if err != nil {
		fmt.Println("Post error:", err)
	}

	// Game state will be updated across all clients using this channel. The clients without priority will block after posting a status update.
	// the client with priority will take thier action and then send another game state update to the rest.
	for {
		fmt.Println(<- resp)
	}

	err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		fmt.Println("write close:", err)
		return
	}
}