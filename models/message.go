package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Username  string `json:"username"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}
