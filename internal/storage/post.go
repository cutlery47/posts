package storage

import (
	"time"

	"github.com/google/uuid"
)

// input-bound Post
type InPost struct {
	// creator id
	UserId uuid.UUID
	// defines if other users can comment on the post
	IsMute bool

	Content string
}

// output-bound Post
// not concurrent-safe by itself!
type Post struct {
	InPost

	Id uuid.UUID

	Upvotes   uint64
	Downvotes uint64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// CommentId -> Comment map
	Comments map[uuid.UUID]Comment
}
