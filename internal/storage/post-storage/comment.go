package storage

import (
	"time"

	"github.com/google/uuid"
)

// input-bound comment
type InComment struct {
	// creator id
	UserId uuid.UUID `json:"user_id"`

	Content string `json:"content"`
}

// output-bound comment
// not concurrent-safe by itself!
type Comment struct {
	InComment

	Id uuid.UUID `json:"id"`

	Upvotes   uint64 `json:"upvotes"`
	Downvotes uint64 `json:"downvotes"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	// CommentId -> Comment
	Replies map[uuid.UUID]Comment `json:"replies"`
}
