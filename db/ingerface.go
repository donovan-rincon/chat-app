package db

import "chat-app/models"

// DBInterface defines the methods that any database implementation must have
type DBInterface interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetOrCreateChatroom(name string) (*models.Chatroom, error)
	GetChatroomByName(name string) (*models.Chatroom, error)
	CreateUserMessage(message *models.UserMessage) error
	GetLastNUserMessages(chatroomID uint, limit int) ([]models.UserMessage, error)
}
