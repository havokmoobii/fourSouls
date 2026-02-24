package gamelogic

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func ClientWelcome() (string, error) {
	fmt.Println("\nPlease enter your username:")
	words := GetInput()
	if len(words) == 0 {
		return "", errors.New("you must enter a username. goodbye")
	}
	username := words[0]
	return username, nil
}

func PrintLobbyHelp() {
	fmt.Println()
	fmt.Println("Possible commands:")
	fmt.Println("* create")
	fmt.Println("* join <room#>")
	fmt.Println("    example:")
	fmt.Println("    join 1")
	fmt.Println("* update")
	fmt.Println("    todo: make lobby update in real time")
	fmt.Println("* quit")
	fmt.Println("* help")
}

func PrintClientHelp() {
	fmt.Println()
	fmt.Println("Possible commands:")
	// Generic game action for test
	fmt.Println("* do")
	fmt.Println("* chat <message>")
	fmt.Println("* dm <user> <message>")
	fmt.Println("* quit")
	fmt.Println("* help")
	fmt.Print("\n> ")
}

func GetInput() []string {
	scanner := bufio.NewScanner(os.Stdin)
	scanned := scanner.Scan()
	if !scanned {
		return nil
	}
	line := scanner.Text()
	line = strings.TrimSpace(line)
	return strings.Fields(line)
	
}