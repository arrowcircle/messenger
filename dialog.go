package main

import (
  "time"
)

type DialogJson struct {
  Id                  uint8         `json:"id"`
  Name                string       `sql:"size:255" json:"name"`
  CreatedAt           time.Time    `json:"created_at"`
  UpdatedAt           time.Time    `json:"updated_at"`
  UserIds             string       `json:"user_ids"`
  LastMessage         string       `json:"last_message"`
  LastMessageID       int          `json:"last_message_id"`
  LastMessageUserID   int          `json:"last_message_user_id"`
  LastSeenMessageID   int          `json:"last_seen_message_id"`
}

type Dialog struct {
  Id              uint8         `json:"id"`
  Name            string       `sql:"size:255" json:"name"`
  CreatedAt       time.Time    `json:"created_at"`
  UpdatedAt       time.Time    `json:"updated_at"`
  Messages        []Message
  LastMessage     Message      `json:"message"`
  LastMessageID   int
}

type DialogUser struct {
  DialogID     uint8
  Dialog       Dialog
  UserID       int
  CreatedAt    time.Time
  UpdatedAt    time.Time
}
