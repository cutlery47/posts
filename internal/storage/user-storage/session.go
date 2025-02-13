package storage

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id     uuid.UUID
	UserId uuid.UUID

	CreatedAt time.Time
	ExpiresAt time.Time
}
