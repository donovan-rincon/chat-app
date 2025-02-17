package db

import (
	"github.com/donovan-rincon/chat-app/database"
	"github.com/donovan-rincon/chat-app/models"
)

// User operations
func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// Chatroom operations
func GetOrCreateChatroom(name string) (*models.Chatroom, error) {
	var chatroom models.Chatroom
	err := database.DB.Where("name = ?", name).FirstOrCreate(&chatroom).Error
	return &chatroom, err
}

func GetChatroomByName(name string) (*models.Chatroom, error) {
	var chatroom models.Chatroom
	err := database.DB.Where("name = ?", name).First(&chatroom).Error
	return &chatroom, err
}

// Message operations
func CreateMessage(message *models.Message) error {
	return database.DB.Create(message).Error
}

func GetLastMessages(chatroomID uint, limit int) ([]models.Message, error) {
	var messages []models.Message
	err := database.DB.Where("chatroom_id = ?", chatroomID).Order("timestamp desc").Limit(limit).Find(&messages).Error
	return messages, err
}
