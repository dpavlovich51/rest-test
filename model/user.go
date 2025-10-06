package model

import (
	"time"
)

type User struct {
	Id        string
	Name      string
	Age       int
	CreatedAt time.Time
}
