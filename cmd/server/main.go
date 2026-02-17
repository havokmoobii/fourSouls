package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
)

type config struct {
	clients map[string]*websocket.Conn
}

func main() {
	cfg := config{
		clients: make(map[string]*websocket.Conn),
	}

	m := http.NewServeMux()

	port := "1337"

	m.HandleFunc("/connect/{username}", cfg.handleConnections)

	srv := http.Server{
		Handler:      m,
		Addr:         ":" + port,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	// this blocks forever, until the server
	// has an unrecoverable error
	fmt.Println("server started on", port)
	err := srv.ListenAndServe()
	log.Fatal(err)
}