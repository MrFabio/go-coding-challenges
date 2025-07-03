package main

import (
	"log"
	"url/api"
)

func main() {
	// Create and start server
	srv := api.NewServer()

	log.Printf("Starting URL Shortener server on port %s", srv.Port)
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
