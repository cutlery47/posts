package mem

import (
	"errors"
	"testing"
	"time"

	storage "github.com/cutlery47/posts/internal/storage/post-storage"
	"github.com/google/uuid"
)

var (
	post storage.Post
)

func TestPostGetEmptyComments(t *testing.T) {
	post = storage.Post{
		Comments: make(map[uuid.UUID]storage.Comment),
	}

	if _, ok := getComment(post, uuid.New()); ok {
		t.Fatalf("should be false")
	}
}

func TestPostUpdateEmptyComments(t *testing.T) {
	post = storage.Post{
		Comments: make(map[uuid.UUID]storage.Comment),
	}

	_, err := updateComment(post, uuid.New(), storage.InComment{})
	if !errors.Is(err, storage.ErrCommNotFound) {
		t.Fatalf("should be %v", storage.ErrCommNotFound)
	}
}

func TestPostDeleteEmptyComments(t *testing.T) {
	post = storage.Post{
		Comments: make(map[uuid.UUID]storage.Comment),
	}

	err := deleteComment(post, uuid.New())
	if !errors.Is(err, storage.ErrCommNotFound) {
		t.Fatalf("should be %v", storage.ErrCommNotFound)
	}
}

func TestPostInsertCommentIntoDeleted(t *testing.T) {
	ts := time.Now()
	post = storage.Post{
		DeletedAt: &ts,
	}

	_, err := insertComment(post, storage.InComment{})
	if !errors.Is(err, storage.ErrPostIsDeleted) {
		t.Fatalf("error: %v", err)
	}
}

func TestPostInsertCommentIntoMute(t *testing.T) {
	post = storage.Post{
		InPost: storage.InPost{
			IsMute: true,
		},
	}

	_, err := insertComment(post, storage.InComment{})
	if !errors.Is(err, storage.ErrPostIsMute) {
		t.Fatalf("error: %v", err)
	}
}

func TestPostDeleteNotFound(t *testing.T) {
	id := uuid.New()

	post = storage.Post{
		Comments: map[uuid.UUID]storage.Comment{
			id: {Id: id},
		},
	}

	err := deleteComment(post, uuid.New())
	if !errors.Is(err, storage.ErrCommNotFound) {
		t.Fatalf("error: %v", err)
	}
}

func TestPostpdateNotFound(t *testing.T) {
	id := uuid.New()

	post = storage.Post{
		Comments: map[uuid.UUID]storage.Comment{
			id: {Id: id},
		},
	}

	_, err := updateComment(post, uuid.New(), storage.InComment{})
	if !errors.Is(err, storage.ErrCommNotFound) {
		t.Fatalf("error: %v", err)
	}
}
