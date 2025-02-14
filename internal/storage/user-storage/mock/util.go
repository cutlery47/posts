package mock

import (
	"context"
	"time"

	storage "github.com/cutlery47/posts/internal/storage/user-storage"
	"github.com/google/uuid"
)

func ctxDone(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return nil
}

func toUser(in storage.InUser) storage.User {
	return storage.User{
		InUser:    in,
		Id:        uuid.New(),
		CreatedAt: time.Now(),
	}
}

func newSession(userId uuid.UUID, expiresAt time.Time) *storage.Session {
	return &storage.Session{
		Id:        uuid.New(),
		UserId:    userId,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}
}
