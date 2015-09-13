package main

import (
  "time"
)

type Message struct {
  Id           int8
  Text         string `sql:"type:text"`
  CreatedAt    time.Time
  UpdatedAt    time.Time
  UserId       int
  DialogId     int8
}
