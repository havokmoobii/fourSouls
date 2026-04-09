package routing

import (
	"fmt"
)

func (cfg *ClientConfig) ReceivePost() {
	// Maybe have a second loop before the game starts for different behavior?
	for {
		var pst Post
		err := cfg.Conn.ReadJSON(&pst)

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
			if cfg.Username == pst.Msg.Sender {
				fmt.Print("\nGameState updated!\n\n> ")
			} else {
				fmt.Print("\n\nGameState updated!\n\n> ")
			}
			cfg.GS = pst.GS
		}

		fmt.Println("\n>")
	}
}
