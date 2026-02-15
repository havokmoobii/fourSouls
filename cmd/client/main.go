package main

// Next time: Make changes so the server can track multiple clients and send a message (and eventually state) to all of them.

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func post(conn *websocket.Conn, msg string) error {
	err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		fmt.Println("Write error:", err)
		return err
	}

	_, reply, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Read error:", err)
		return err
	}

	fmt.Println(string(reply))
	return nil
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:1337/ws", nil)
	if err != nil {
		fmt.Println("Dial error:", err)
		return
	}
	defer conn.Close()
	
	msg := "Hey does this work???"

	err = post(conn, msg)
	if err != nil {
		fmt.Println("Post error:", err)
	}

}