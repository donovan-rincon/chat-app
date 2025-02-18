package main

import (
	"log"

	"chat-app/database"
	"chat-app/db"
	"chat-app/handlers"
)

func main() {
	database.Init()
	dbInstance := db.NewGormDB()

	// Setup routes
	r := handlers.SetupRouter(dbInstance)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
