package storage

import (
	"time"

	"github.com/google/uuid"
)

// input-bound Post
type InPost struct {
	// creator id
	UserId uuid.UUID `json:"user_id"`
	// defines if other users can comment on the post
	IsMute bool `json:"is_mute"`

	Content string `json:"content"`
}

// output-bound Post
// not concurrent-safe by itself!
type Post struct {
	InPost `json:"in_post"`

	Id uuid.UUID `json:"id"`

	Upvotes   uint64 `json:"upvotes"`
	Downvotes uint64 `json:"downvotes"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	// CommentId -> Comment map
	Comments map[uuid.UUID]Comment `json:"comments"`
}
