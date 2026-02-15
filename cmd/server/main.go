package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	m := http.NewServeMux()

	port := "1337"

	m.HandleFunc("/ws", handlerWebsocket)

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