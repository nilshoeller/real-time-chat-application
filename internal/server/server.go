package server

import (
	"log"
	"net/http"

	"github.com/nilshoeller/real-time-chat-application/internal/routes"
)


func Run() {
	// Create a new router
	mux := http.NewServeMux()

	// Serve static files from the static directory
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Routing
	mux.HandleFunc("/", routes.RouteIndexPage())

	// Start the server
	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}