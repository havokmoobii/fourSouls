package main

// Next Time:
// Update server text to be less spammy.
// Fix crash if someone leaves a started game and attempts to rejoin.
// Add status to rooms.
// Do same text formatting for game loop as for the lobby loop
// Make it so that the game doesn't start until all players have passed
// the "Press enter to continue message"?
// Maybe have a message post saying x player is ready when the press enter
// and after the last player does the console will say game is starting!

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
		StartSignal: false,
		CloseSignal: false,
	}

	fmt.Println("Welcome to the Four Souls client!")

	var err error
	cfg.Username, err = gamelogic.ClientWelcome()
	if err != nil {
		return
	}

	fmt.Println()
	gamelogic.PrintLobbyHelp()

	err = cfg.CheckServer()
	if err != nil {
		fmt.Print("\nerror: ", err)
		fmt.Println()
	}

	// Lobby loop
	for {
		fmt.Println()
		fmt.Print("> ")
		words := gamelogic.GetInput()
		if cfg.StartSignal {
			break
		}
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "create":
			if cfg.RoomNumber != 0 {
				fmt.Println("Cannot join multiple rooms!")
				continue
			}
			err = cfg.CreateRoom()
			if err != nil {
				fmt.Print("\nerror: ", err)
				fmt.Println()
			}
		case "join":
			if len(words) < 2 {
				fmt.Println("Join command must include a room number!")
				continue
			}
			if cfg.RoomNumber != 0 {
				fmt.Println("Cannot join multiple rooms!")
				continue
			}
			err = cfg.JoinRoom(words[1])
			if err != nil {
				fmt.Print("\nerror: ", err)
				fmt.Println()
			}
		case "update":
			err = cfg.CheckServer()
			if err != nil {
				fmt.Print("\nerror: ", err)
				fmt.Println()
			}
		case "start":
			err = cfg.SendPost(routing.Post{
				Kind: routing.PostGameStart,
			})
			if err != nil {
				fmt.Print("error: ", err)
			}
		case "quit":
			return
		case "help":
			gamelogic.PrintLobbyHelp()
			fmt.Println()
		default:
			fmt.Println("\nUnknown command")
		}
	}
	defer cfg.Conn.Close()

	fmt.Println()
	gamelogic.PrintClientHelp()

	// Game Loop
	for {
		fmt.Println()
		fmt.Print("> ")
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "do":
			err = cfg.SendPost(routing.Post{
				Kind: routing.PostStateUpdate,
				GS:   cfg.GS,
				Msg: routing.Message{
					Sender: cfg.Username,
				},
			})
			if err != nil {
				fmt.Print("error: ", err)
			}
		case "chat":
			if len(words) > 1 {
				err = cfg.SendPost(routing.Post{
					Kind: routing.PostChat,
					Msg: routing.Message{
						Sender: cfg.Username,
						Body:   strings.Join(words[1:], " "),
					},
				})
				if err != nil {
					fmt.Print("error: ", err)
				}
			} else {
				fmt.Print("\nerror: 'chat' must be followed by a message!\n\n> ")
			}
		case "dm":
			if len(words) > 2 {
				// Once usernames are tracked in gamestate, check if valid recipient here
				err = cfg.SendPost(routing.Post{
					Kind: routing.PostChat,
					Msg: routing.Message{
						Sender:    cfg.Username,
						Recipient: words[1],
						Body:      strings.Join(words[2:], " "),
					},
				})
				if err != nil {
					fmt.Print("error: ", err)
				}
			} else {
				fmt.Print("\nerror: 'dm' must be followed by a username and a message!\n\n> ")
			}
		case "quit":
			cfg.CloseSignal = true
			err = cfg.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("write close:", err)
			}
			return
		case "help":
			gamelogic.PrintClientHelp()
		default:
			fmt.Println("Unknown command")
		}
	}
}
