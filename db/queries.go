package db

import (
	"chat-app/database"
	"chat-app/models"

	"gorm.io/gorm"
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
	err := database.DB.Where("name = ?", name).First(&chatroom).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if chatroom.ID == 0 {
		chatroom.Name = name
		err = database.DB.Create(&chatroom).Error
		if err != nil {
			return nil, err
		}
	}
	return &chatroom, nil
}

func GetChatroomByName(name string) (*models.Chatroom, error) {
	var chatroom models.Chatroom
	err := database.DB.Where("name = ?", name).First(&chatroom).Error
	return &chatroom, err
}

// Message operations
func CreateUserMessage(message *models.UserMessage) error {
	return database.DB.Create(message).Error
}

func GetLastNUserMessages(chatroomID uint, limit int) ([]models.UserMessage, error) {
	var messages []models.UserMessage
	err := database.DB.Where("chatroom_id = ?", chatroomID).Order("timestamp desc").Limit(limit).Find(&messages).Error
	return messages, err
}
