package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type apiConfig struct {
	test string
}

func main() {
	m := http.NewServeMux()

	cfg := apiConfig{}

	port := "1337"

	m.HandleFunc("GET /test", cfg.handlerTest)

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

func (cfg *apiConfig) handlerTest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello World!")
}