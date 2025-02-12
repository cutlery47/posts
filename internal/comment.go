package internal

import (
	"time"

	"github.com/google/uuid"
)

// input-bound comment
type InComment struct {
	// creator id
	userId uuid.UUID

	content string
}

// output-bound comment
// not concurrent-safe by itself!
type Comment struct {
	InComment

	id uuid.UUID

	upvotes   uint64
	downvotes uint64

	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time

	// CommentId -> Comment
	replies map[uuid.UUID]Comment
}

func (c Comment) insertReply(in InComment) Comment {
	var (
		comm = toComment(in)
	)

	// loop until no collisions detected
	for _, ok := c.replies[comm.id]; ok; _, ok = c.replies[comm.id] {
		comm = toComment(in)
	}

	return comm
}

func toComment(in InComment) Comment {
	return Comment{
		id:        uuid.New(),
		upvotes:   0,
		downvotes: 0,
		createdAt: time.Now(),
		updatedAt: time.Now(),
		deletedAt: nil,
		replies:   make(map[uuid.UUID]Comment),
		InComment: in,
	}
}
