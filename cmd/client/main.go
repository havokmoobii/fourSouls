package main

// Next time: Set up a repl loop and give the clients options.
// I want to set up chat, but not sure how with the channel blocking implementation

import (
	"fmt"
	"strings"
	"github.com/gorilla/websocket"
	"github.com/havokmoobii/fourSouls/internal/gamelogic"
)

func main() {
	cfg := clientConfig{
		gs: gamelogic.GameState{},
	}
	
	err := cfg.connect()
	if err != nil {
		// Error is handled in the function, so we can simply return.
		return
	}
	defer cfg.conn.Close()

	fmt.Println("Success!")

	go cfg.receivePost()
	go cfg.receiveChatPost()

	// When player has priorty they will end each action with a call to post to update the rest of the players and pass priorty.
	err = cfg.post(cfg.gs)
	if err != nil {
		fmt.Println("Post error:", err)
	}

	// Game state will be updated after each priority player action. Below will be a REPL loop. Messages to the clients will always provide a cursor to
	// to make the experience mostly seamless. The idea is that the player is told it is their turn and they can enter relevant commands, but can still use
	// chat and some other features when they do not have priority
	for {
		msg := gamelogic.GetInput()
		cfg.chatPost(strings.Join(msg, " "))
		msg = gamelogic.GetInput()
		if len(msg) > 1 {
			cfg.chatDM(msg[0], strings.Join(msg[1:], " "))
		}
	}

	err = cfg.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		fmt.Println("write close:", err)
		return
	}
}