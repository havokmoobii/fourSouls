package main

// Next time: Set up making a game room. Limit players to 4. Maybe make server able to host multiple games.

import (
	"fmt"
	"strings"
	"github.com/gorilla/websocket"
	"github.com/havokmoobii/fourSouls/internal/gamelogic"
	"github.com/havokmoobii/fourSouls/internal/routing"
)

func main() {
	cfg := routing.ClientConfig{
		GS:          gamelogic.GameState{},
		CloseSignal: false,
	}

	fmt.Println("Welcome to the Four Souls client!")

	cfg.CheckServer()
	
	err := cfg.Connect()
	if err != nil {
		// Error is handled in the function, so we can simply return.
		return
	}
	defer cfg.Conn.Close()

	fmt.Println("Success!")

	go cfg.ReceivePost()
	go cfg.ReceiveChatPost()

	// When player has priorty they will end each action with a call to post to update the rest of the players and pass priorty.
	// Game state will be updated after each priority player action. Below will be a REPL loop. Messages to the clients will always provide a cursor to
	// to make the experience mostly seamless. The idea is that the player is told it is their turn and they can enter relevant commands, but can still use
	// chat and some other features when they do not have priority
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "chat":
			if len(words) > 1{
				err = cfg.ChatPost(strings.Join(words[1:], " "))
				if err != nil {
					fmt.Print("error:", err)
				}
			} 
		case "quit":
			cfg.CloseSignal = true
			err = cfg.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("write close:", err)
			}
			err = cfg.ChatConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("write close:", err)
			}
			return
		default:
			fmt.Println("Unknown command")
		}
	}
}