package main

import (
	"log"

	"chat-app/database"
	"chat-app/handlers"
)

func main() {
	database.Init()

	// Setup routes
	r := handlers.SetupRouter()

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
