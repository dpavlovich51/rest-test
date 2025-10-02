package model

import (
	"time"
	"github.com/google/uuid"
)

type Message struct {
	Id       uuid.UUID
	UserId   string
	Body     string
	CreateAt time.Time
}
