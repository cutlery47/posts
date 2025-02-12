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

func (p Post) getComment(id uuid.UUID) (*Comment, bool) {
	if len(p.comments) == 0 {
		return nil, false
	}

	var (
		comm *Comment
	)

	// iterate down the tree for each comment
	for idx, c := range p.comments {
		if idx == id {
			comm = &c
			break
		}

		r, ok := c.getReply(id)
		if ok {
			comm = r
			break
		}
	}

	if comm == nil {
		return nil, false
	}

	return comm, true
}

func (p Post) updateComment(id uuid.UUID, in InComment) (*Comment, error) {
	if len(p.comments) == 0 {
		return nil, ErrCommNotFound
	}

	var (
		comm *Comment
	)

	// iterate down the tree and update
	for idx, c := range p.comments {
		if idx == id {
			c.content = in.content
			c.updatedAt = time.Now()
			p.comments[idx] = c
			comm = &c
			break
		}

		r, err := c.updateReply(id, in)
		if err == nil {
			comm = r
			break
		}
	}

	if comm == nil {
		return nil, ErrCommNotFound
	}

	return comm, nil
}

func (p Post) deleteComment(id uuid.UUID) error {
	if len(p.comments) == 0 {
		return ErrCommNotFound
	}

	for idx, c := range p.comments {
		if idx == id {
			ts := time.Now()
			c.deletedAt = &ts
			p.comments[idx] = c
			return nil
		}

		err := c.deleteReply(id)
		if err == nil {
			return nil
		}
	}

	return ErrCommNotFound
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
