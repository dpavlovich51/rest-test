package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time
}
