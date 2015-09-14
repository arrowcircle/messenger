package main

import (
  "time"
)

type Message struct {
  Id           int8           `json:"id"`
  Text         string         `sql:"type:text" json:"text"`
  CreatedAt    time.Time      `json:"created_at"`
  UpdatedAt    time.Time      `json:"updated_at"`
  UserID       int            `json:"user_id"`
  DialogID     int8           `json:"dialog_id"`
  Dialog       Dialog
}
