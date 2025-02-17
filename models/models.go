package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserMessage struct {
	gorm.Model
	ChatroomID uint   `json:"chatroom_id"`
	UserID     uint   `json:"user_id"`
	Username   string `json:"username"`
	Message    string `json:"message"`
	Timestamp  string `json:"timestamp"`
}

type Chatroom struct {
	gorm.Model
	Name string `json:"name"`
}

type WSMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}
