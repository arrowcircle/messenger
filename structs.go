package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Impl used to provide handler to DB
type Impl struct {
	DB *gorm.DB
}

// UserJSON is used for empty requests
type UserJSON struct {
	ID           int `json:"id"`
	DialogsCount int `json:"dialogs_count"`
}

// DialogJSON used for index action of API
type DialogJSON struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	UserIds           string    `json:"user_ids"`
	LastMessage       string    `json:"last_message"`
	LastMessageID     int       `json:"last_message_id"`
	LastMessageUserID int       `json:"last_message_user_id"`
	LastSeenMessageID int       `json:"last_seen_message_id"`
}

// DialogCreateJSON is used for dialogs creation
type DialogCreateJSON struct {
	Name    string `json:"name"`
	UserIds []int  `json:"user_ids"`
	Message string `json:"message"`
}

// Dialog is used to save dialogs into DB via GORM
type Dialog struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	LastMessageID int       `json:"last_message_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// DialogShowJSON is used to form json
type DialogShowJSON struct {
	ID       int           `json:"id"`
	Name     string        `json:"name"`
	Messages []MessageJSON `json:"messages"`
}

// MessageJSON is used to response message in JSON format
type MessageJSON struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int       `json:"user_id"`
}

// Message is used to put messages in DB via GORM
type Message struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	UserID   int    `json:"user_id"`
	DialogID int    `json:"dialog_id"`
}
