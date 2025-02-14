package storage

import (
	"time"

	"github.com/google/uuid"
)

type InUser struct {
	Name string `json:"name"` // unique
	Role string `json:"role"`
}

type User struct {
	InUser `json:"in_user"`

	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	AdminRole = "admin"
	UserRole  = "user"
)
