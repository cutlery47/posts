package mem

import (
	"time"

	storage "github.com/cutlery47/posts/internal/storage/post-storage"
	"github.com/google/uuid"
)

func insertReply(c storage.Comment, in storage.InComment) (*storage.Comment, error) {
	if c.DeletedAt != nil {
		return nil, storage.ErrCommIsDeleted
	}

	var (
		repl storage.Comment = toComment(in)
	)

	// loop until no collisions detected
	for _, ok := c.Replies[repl.Id]; ok; _, ok = c.Replies[repl.Id] {
		repl = toComment(in)
	}

	c.Replies[repl.Id] = repl

	return &repl, nil
}

func getReply(c storage.Comment, id uuid.UUID) (*storage.Comment, bool) {
	if len(c.Replies) == 0 {
		return nil, false
	}

	var (
		repl *storage.Comment
	)

	// iterate down the tree for each reply
	for idx, r := range c.Replies {
		if idx == id {
			repl = &r
			break
		}

		cr, ok := getReply(r, id)
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

func updateReply(c storage.Comment, id uuid.UUID, in storage.InComment) (*storage.Comment, error) {
	if len(c.Replies) == 0 {
		return nil, storage.ErrCommNotFound
	}

	var (
		repl *storage.Comment
	)

	// iterate down the tree and update
	for idx, r := range c.Replies {
		if idx == id {
			r.Content = in.Content
			r.UpdatedAt = time.Now()
			c.Replies[idx] = r
			repl = &r
			break
		}

		rr, err := updateReply(r, id, in)
		if err == nil {
			repl = rr
			break
		}
	}

	if repl == nil {
		return nil, storage.ErrCommNotFound
	}

	return repl, nil
}

func deleteReply(c storage.Comment, id uuid.UUID) error {
	if len(c.Replies) == 0 {
		return storage.ErrCommNotFound
	}

	for idx, r := range c.Replies {
		if idx == id {
			ts := time.Now()
			r.DeletedAt = &ts
			r.Replies[idx] = r
			return nil
		}

		err := deleteReply(r, id)
		if err == nil {
			return nil
		}
	}

	return storage.ErrCommNotFound
}
