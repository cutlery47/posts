package mem

import (
	"time"

	storage "github.com/cutlery47/posts/internal/storage/post-storage"
	"github.com/google/uuid"
)

func insertComment(p storage.Post, in storage.InComment) (*storage.Comment, error) {
	if p.DeletedAt != nil {
		return nil, storage.ErrPostIsDeleted
	}

	if p.IsMute {
		return nil, storage.ErrPostIsMute
	}

	var (
		comm storage.Comment = toComment(in)
	)

	// loop until no collisions detected
	for _, ok := p.Comments[comm.Id]; ok; _, ok = p.Comments[comm.Id] {
		comm = toComment(in)
	}

	p.Comments[comm.Id] = comm

	return &comm, nil
}

func getComment(p storage.Post, id uuid.UUID) (*storage.Comment, bool) {
	if len(p.Comments) == 0 {
		return nil, false
	}

	var (
		comm *storage.Comment
	)

	// iterate down the tree for each comment
	for idx, c := range p.Comments {
		if idx == id {
			comm = &c
			break
		}

		r, ok := getReply(c, id)
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

func updateComment(p storage.Post, id uuid.UUID, in storage.InComment) (*storage.Comment, error) {
	if len(p.Comments) == 0 {
		return nil, storage.ErrCommNotFound
	}

	var (
		comm *storage.Comment
	)

	// iterate down the tree and update
	for idx, c := range p.Comments {
		if idx == id {
			c.Content = in.Content
			c.UpdatedAt = time.Now()
			p.Comments[idx] = c
			comm = &c
			break
		}

		r, err := updateReply(c, id, in)
		if err == nil {
			comm = r
			break
		}
	}

	if comm == nil {
		return nil, storage.ErrCommNotFound
	}

	return comm, nil
}

func deleteComment(p storage.Post, id uuid.UUID) error {
	if len(p.Comments) == 0 {
		return storage.ErrCommNotFound
	}

	for idx, c := range p.Comments {
		if idx == id {
			ts := time.Now()
			c.DeletedAt = &ts
			p.Comments[idx] = c
			return nil
		}

		err := deleteReply(c, id)
		if err == nil {
			return nil
		}
	}

	return storage.ErrCommNotFound
}
