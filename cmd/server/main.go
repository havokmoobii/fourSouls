package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
	"github.com/havokmoobii/fourSouls/internal/routing"
)

func main() {
	cfg := routing.ServerConfig{
		Clients:     make(map[string]*websocket.Conn),
		ChatClients: make(map[string]*websocket.Conn),
	}

	m := http.NewServeMux()

	port := "1337"

	m.HandleFunc("/status", cfg.HandleStatus)
	m.HandleFunc("/connect/{username}", cfg.HandleConnections)
	m.HandleFunc("/chat/connect/{username}", cfg.HandleChatConnections)

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