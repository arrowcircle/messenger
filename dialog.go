package main

import (
  "time"
)

type Dialog struct {
  Id           int8
  Name         string  `sql:"size:255"`
  CreatedAt    time.Time
  UpdatedAt    time.Time
  Messages     []Message
}

type DialogUser struct {
  DialogId     int8
  UserId       int
  CreatedAt    time.Time
  UpdatedAt    time.Time
}
