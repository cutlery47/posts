package storage

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id     uuid.UUID `json:"id"`
	UserId uuid.UUID `json:"user_id"`

	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
