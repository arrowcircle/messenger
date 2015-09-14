package main

import (
  "time"
)

type Dialog struct {
  Id           int8         `json:"id"`
  Name         string       `sql:"size:255" json:"name"`
  CreatedAt    time.Time    `json:"created_at"`
  UpdatedAt    time.Time    `json:"updated_at"`
  Messages     []Message
}

type DialogUser struct {
  DialogID     int8
  Dialog       Dialog
  UserID       int
  CreatedAt    time.Time
  UpdatedAt    time.Time
}
