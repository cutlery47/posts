package storage

import (
	"time"

	"github.com/google/uuid"
)

// input-bound comment
type InComment struct {
	// creator id
	UserId uuid.UUID

	Content string
}

// output-bound comment
// not concurrent-safe by itself!
type Comment struct {
	InComment

	Id uuid.UUID

	Upvotes   uint64
	Downvotes uint64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// CommentId -> Comment
	Replies map[uuid.UUID]Comment
}
