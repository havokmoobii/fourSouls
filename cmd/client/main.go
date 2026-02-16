package main

// Next time: Make changes so the server can track multiple clients and send a message to all of them.

import (
	"fmt"
	"os"
	"bufio"
	"github.com/gorilla/websocket"
	"github.com/havokmoobii/fourSouls/internal/gamelogic"
)

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
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:1337/ws", nil)
	if err != nil {
		fmt.Println("Dial error:", err)
		return
	}
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)

	gs := gamelogic.GameState{
		Player: "HavokMoobii",
	}

	fmt.Println(gs)

	fmt.Println("Press enter to send a message to the server.")
	scanner.Scan()

	resp, err := post(conn, gs)
	if err != nil {
		fmt.Println("Post error:", err)
	}

	fmt.Println(resp)

	err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		fmt.Println("write close:", err)
		return
	}

	

}