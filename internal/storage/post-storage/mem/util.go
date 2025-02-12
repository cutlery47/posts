package mem

import (
	"context"
	"time"

	storage "github.com/cutlery47/posts/internal/storage/post-storage"
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

func toComment(in storage.InComment) storage.Comment {
	return storage.Comment{
		Id:        uuid.New(),
		Upvotes:   0,
		Downvotes: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
		Replies:   make(map[uuid.UUID]storage.Comment),
		InComment: in,
	}
}

func toPost(in storage.InPost) storage.Post {
	return storage.Post{
		Id:        uuid.New(),
		Upvotes:   0,
		Downvotes: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
		Comments:  make(map[uuid.UUID]storage.Comment),
		InPost:    in,
	}
}
