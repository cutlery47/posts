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
		repl = toComment(in)
	)

	// loop until no collisions detected
	for _, ok := c.replies[repl.id]; ok; _, ok = c.replies[repl.id] {
		repl = toComment(in)
	}

	return repl
}

func (c Comment) getReply(id uuid.UUID) (*Comment, bool) {
	if len(c.replies) == 0 {
		return nil, false
	}

	var (
		repl *Comment
	)

	// iterate down the tree for each reply
	for idx, r := range c.replies {
		if idx == id {
			repl = &r
			break
		}

		cr, ok := r.getReply(id)
		if ok {
			repl = cr
			break
		}
	}

	if repl == nil {
		return nil, false
	}

	return repl, true
}

func (c Comment) updateReply(id uuid.UUID, in InComment) (*Comment, error) {
	if len(c.replies) == 0 {
		return nil, ErrCommNotFound
	}

	var (
		repl *Comment
	)

	// iterate down the tree and update
	for idx, r := range c.replies {
		if idx == id {
			r.content = in.content
			r.updatedAt = time.Now()
			c.replies[idx] = r
			repl = &r
			break
		}

		rr, err := r.updateReply(id, in)
		if err != nil {
			repl = rr
			break
		}
	}

	if repl == nil {
		return nil, ErrCommNotFound
	}

	return repl, nil
}

func (c Comment) deleteReply(id uuid.UUID) error {
	if len(c.replies) == 0 {
		return ErrCommNotFound
	}

	for idx, r := range c.replies {
		if idx == id {
			ts := time.Now()
			r.deletedAt = &ts
			r.replies[idx] = r
			return nil
		}

		err := r.deleteReply(id)
		if err == nil {
			return nil
		}
	}

	return ErrCommNotFound
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
