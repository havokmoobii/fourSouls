package main

// But Actually Next Time: Implement ReceiveCombinedPost
// Next time: Figure out timing with starting game and the menu. Currently it loops back to the menu before the start game command registers
// Idea: Use a channel for the above problem. Have the program halt after the start command and the reciever can send the proceed command. Would probably cause
//     Issues with the other clients though.
// Idea: have a game started flag in client config to tell the WS listener to change its behavior after leaving the lobby to reduce the number of connections

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

	gamelogic.PrintLobbyHelp()
	err = cfg.CheckServer()
	if err != nil {
		fmt.Print("\nerror: ", err)
		fmt.Println("\n")
	}

	for {
		fmt.Print("> ")
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "create":
			err = cfg.CreateRoom()
			if err != nil {
				fmt.Print("\nerror: ", err)
				fmt.Println("\n")
			}
		case "join":
			err = cfg.Connect()
			if err != nil {
				// Error is handled in the function, so we can simply return.
				return
			}
			go cfg.ReceivePost()
		case "update":
			err = cfg.CheckServer()
			if err != nil {
				fmt.Print("\nerror: ", err)
				fmt.Println("\n")
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
			fmt.Println("Unknown command")
		}
		if cfg.StartSignal {
			break
		}
	} 
	defer cfg.Conn.Close()

	
	gamelogic.PrintClientHelp()

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
		case "do":
			err = cfg.SendPost(routing.Post{
				Kind: routing.PostStateUpdate,
				GS:   cfg.GS,
				Msg:  routing.Message{
					Sender: cfg.Username,
				},
			})
			if err != nil {
					fmt.Print("error: ", err)
				}
		case "chat":
			if len(words) > 1{
				err = cfg.SendPost(routing.Post{
					Kind:   routing.PostChat,
					Msg:    routing.Message{
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
			if len(words) > 2{
				// Once usernames are tracked in gamestate, check if valid recipient here
				err = cfg.SendPost(routing.Post{
					Kind:   routing.PostChat,
					Msg:    routing.Message{
						Sender: cfg.Username,
						Recipient: words[1],
						Body:   strings.Join(words[2:], " "),
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