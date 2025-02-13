package storage

import (
	"time"

	"github.com/google/uuid"
)

type InUser struct {
	Role string
}

type User struct {
	InUser

	Id uuid.UUID

	CreatedAt time.Time
}

var (
	AdminRole = "admin"
	UserRole  = "user"
)
