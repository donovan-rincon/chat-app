package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

type Message struct {
	gorm.Model
	ChatroomID uint   `json:"chatroom_id"`
	UserID  string `json:"userid"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type Chatroom struct {
	gorm.Model
	Name string `json:"name"`
}
