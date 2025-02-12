package internal

import (
	"time"

	"github.com/google/uuid"
)

// input-bound Post
type InPost struct {
	// creator id
	userId uuid.UUID
	// defines if other users can comment on the post
	isMute bool

	content string
}

// output-bound Post
// not concurrent-safe by itself!
type Post struct {
	InPost

	id uuid.UUID

	upvotes   uint64
	downvotes uint64

	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time

	// CommentId -> Comment map
	comments map[uuid.UUID]Comment
}

func (p Post) insertComment(in InComment) Comment {
	var (
		comm = toComment(in)
	)

	// loop until no collisions detected
	for _, ok := p.comments[comm.id]; ok; _, ok = p.comments[comm.id] {
		comm = toComment(in)
	}

	return comm
}

func toPost(in InPost) Post {
	return Post{
		id:        uuid.New(),
		upvotes:   0,
		downvotes: 0,
		createdAt: time.Now(),
		updatedAt: time.Now(),
		deletedAt: nil,
		comments:  make(map[uuid.UUID]Comment),
		InPost:    in,
	}
}
