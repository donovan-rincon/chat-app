package main

import (
	"log"

	"github.com/donovan-rincon/chat-app/database"
	"github.com/donovan-rincon/chat-app/handlers"
)

func main() {
	database.Init()

	// Setup routes
	r := handlers.SetupRouter()

	// start the bot
	// bot.BotListener()

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
