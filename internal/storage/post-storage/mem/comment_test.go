package mem

import (
	"errors"
	"testing"
	"time"

	storage "github.com/cutlery47/posts/internal/storage/post-storage"
	"github.com/google/uuid"
)

var (
	comm storage.Comment
)

func TestCommentGetEmptyReplies(t *testing.T) {
	ts := time.Now()
	comm = storage.Comment{
		DeletedAt: &ts,
	}

	_, err := insertReply(comm, storage.InComment{})
	if !errors.Is(err, storage.ErrCommIsDeleted) {
		t.Fatalf("error: %v", err)
	}
}

func TestCommentGetReply(t *testing.T) {
	id := uuid.New()

	comm = storage.Comment{
		Replies: map[uuid.UUID]storage.Comment{
			id: storage.Comment{Id: id},
		},
	}

	_, ok := getReply(comm, id)
	if !ok {
		t.Fatal("not found")
	}
}

func TestCommentNestedGetReply(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()

	comm = storage.Comment{
		Replies: map[uuid.UUID]storage.Comment{
			id1: storage.Comment{
				Id: id1,
				Replies: map[uuid.UUID]storage.Comment{
					id2: storage.Comment{
						Id: id2,
					},
				},
			},
		},
	}

	_, ok := getReply(comm, id2)
	if !ok {
		t.Fatal("not found")
	}
}

func TestCommentNestedUpdateReply(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()

	comm = storage.Comment{
		Replies: map[uuid.UUID]storage.Comment{
			id1: storage.Comment{
				Id: id1,
				Replies: map[uuid.UUID]storage.Comment{
					id2: storage.Comment{
						Id: id2,
					},
				},
			},
		},
	}

	_, err := updateReply(comm, id2, storage.InComment{})
	if err != nil {
		t.Fatal("error: ", err)
	}
}

func TestCommentNestedUpdateNonexistantReply(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()

	comm = storage.Comment{
		Replies: map[uuid.UUID]storage.Comment{
			id1: storage.Comment{
				Id: id1,
				Replies: map[uuid.UUID]storage.Comment{
					id2: storage.Comment{
						Id: id2,
					},
				},
			},
		},
	}

	_, err := updateReply(comm, uuid.New(), storage.InComment{})
	if !errors.Is(err, storage.ErrCommNotFound) {
		t.Fatal("error: ", err)
	}
}

func TestCommentNestedDeleteReply(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()

	comm = storage.Comment{
		Replies: map[uuid.UUID]storage.Comment{
			id1: storage.Comment{
				Id: id1,
				Replies: map[uuid.UUID]storage.Comment{
					id2: storage.Comment{
						Id:      id2,
						Replies: map[uuid.UUID]storage.Comment{},
					},
				},
			},
		},
	}

	err := deleteReply(comm, id2)
	if err != nil {
		t.Fatal("error: ", err)
	}
}

func TestCommentNestedDeleteNonexistantReply(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()

	comm = storage.Comment{
		Replies: map[uuid.UUID]storage.Comment{
			id1: storage.Comment{
				Id: id1,
				Replies: map[uuid.UUID]storage.Comment{
					id2: storage.Comment{
						Id: id2,
					},
				},
			},
		},
	}

	err := deleteReply(comm, uuid.New())
	if !errors.Is(err, storage.ErrCommNotFound) {
		t.Fatal("error: ", err)
	}
}
