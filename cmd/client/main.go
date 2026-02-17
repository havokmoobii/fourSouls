package main

// Next time: Make changes so the server can track multiple clients and send a message to all of them.

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

func post(conn *websocket.Conn, msg interface{}) (interface{}, error) {
	err := conn.WriteJSON(msg)
	if err != nil {
		fmt.Println("Write error:", err)
		return nil, err
	}

	var resp interface{}
	err = conn.ReadJSON(&resp)
	if err != nil {
		fmt.Println("Read error:", err)
		return nil, err
	}

	return resp, nil
}

func main() {
	conn, err := connect()
	if err != nil {
		// Error is handled in the function, so we can simply return.
		return
	}
	defer conn.Close()

	fmt.Println("Success!")

	gs := gamelogic.GameState{
		Player: "HavokMoobii",
	}

	resp, err := post(conn, gs)
	if err != nil {
		fmt.Println("Post error:", err)
	}

	fmt.Println(resp)

	for {}

	err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		fmt.Println("write close:", err)
		return
	}

	

}