package storage

import (
	"context"

	"github.com/google/uuid"
)

type Storage interface {
	// registeres user with given inputs
	Register(ctx context.Context, in InUser) (*User, error)
	// logs given user in
	Login(ctx context.Context, in InUser) (*Session, error)
	// logs given user out
	Logout(ctx context.Context, sesh Session) error

	GetSession(ctx context.Context, id uuid.UUID) (*Session, error)
}
