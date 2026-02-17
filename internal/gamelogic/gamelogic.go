package gamelogic

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func ClientWelcome() (string, error) {
	fmt.Println("Welcome to the Four Souls client!")
	fmt.Println("Please enter your username:")
	words := GetInput()
	if len(words) == 0 {
		return "", errors.New("you must enter a username. goodbye")
	}
	username := words[0]
	return username, nil
}

func GetInput() []string {
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanned := scanner.Scan()
	if !scanned {
		return nil
	}
	line := scanner.Text()
	line = strings.TrimSpace(line)
	return strings.Fields(line)
}