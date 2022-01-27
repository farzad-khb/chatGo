package main

import (
	"log"
	"net/http"
	"os"
	"ws/internal/handlers"
)

func main() {
	mux := routes()
	port, ok := os.LookupEnv("PORT")
	// default port
	if !ok {
		port = "8080"

	}

	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()

	// starting a server
	log.Println("Starting web server on port :", port)
	_ = http.ListenAndServe(":"+port, mux)

}
