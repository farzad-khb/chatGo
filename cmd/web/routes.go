package main

import (
	"net/http"
	"ws/internal/handlers"

	"github.com/bmizerany/pat"
)

// define our routes here
func routes() http.Handler {

	// multiplexer
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndpoint))

	return mux
}
