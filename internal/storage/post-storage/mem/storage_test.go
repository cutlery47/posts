package mem_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/cutlery47/posts/config"
	storage "github.com/cutlery47/posts/internal/storage/post-storage"
	"github.com/cutlery47/posts/internal/storage/post-storage/mem"
	"github.com/google/uuid"
)

var (
	store storage.Storage

	conf = config.PostStorage{
		DumpEnabled: false,
	}

	w io.Writer
	r io.Reader
)

func TestStorageInsertPost(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	in := storage.InPost{
		UserId:  uuid.New(),
		IsMute:  false,
		Content: "content",
	}

	_, err := store.InsertPost(ctx, in)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestStorageGetPost(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	in := storage.InPost{
		UserId:  uuid.New(),
		IsMute:  false,
		Content: "content",
	}

	post, err := store.InsertPost(ctx, in)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	gPost, err := store.GetPost(ctx, post.Id)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if gPost.InPost != in {
		t.Fatalf("wrong post")
	}
}

func TestStorageGetNonexistantPost(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	_, err := store.InsertPost(ctx, storage.InPost{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = store.GetPost(ctx, uuid.New())
	if !(errors.Is(err, storage.ErrPostNotFound)) {
		t.Fatalf("error: %v", err)
	}
}

func TestStorageGetPosts(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	in1 := storage.InPost{
		UserId:  uuid.New(),
		IsMute:  false,
		Content: "content1",
	}

	in2 := storage.InPost{
		UserId:  uuid.New(),
		IsMute:  true,
		Content: "content2",
	}

	post1, err := store.InsertPost(ctx, in1)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	post2, err := store.InsertPost(ctx, in2)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	posts, err := store.GetPosts(ctx)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if len(posts) != 2 {
		t.Fatalf("wrong length")
	}

	if post1.InPost != in1 || post2.InPost != in2 {
		t.Fatalf("wrong posts")
	}
}

func TestStorageDeletePost(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	post, err := store.InsertPost(ctx, storage.InPost{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = store.DeletePost(ctx, post.Id)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	del, err := store.GetPost(ctx, post.Id)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if del.DeletedAt == nil {
		t.Fatalf("deletion didn't persist")
	}
}

func TestStorageDeleteNonexistantPost(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	_, err := store.InsertPost(ctx, storage.InPost{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = store.DeletePost(ctx, uuid.New())
	if !errors.Is(err, storage.ErrPostNotFound) {
		t.Fatalf("error: %v", err)
	}
}

func TestStorageDeleteDeletedPost(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	post, err := store.InsertPost(ctx, storage.InPost{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = store.DeletePost(ctx, post.Id)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = store.DeletePost(ctx, post.Id)
	if !errors.Is(err, storage.ErrPostIsDeleted) {
		t.Fatalf("managed to delete deleted post: %v", err)
	}
}

func TestStorageUpdatePost(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	id := uuid.New()

	in := storage.InPost{
		UserId:  id,
		IsMute:  false,
		Content: "content",
	}

	post, err := store.InsertPost(ctx, in)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	inUpd := storage.InPost{
		UserId:  id,
		IsMute:  true,
		Content: "skibidi",
	}

	upd, err := store.UpdatePost(ctx, post.Id, inUpd)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if upd.InPost != inUpd {
		t.Fatalf("updates didn't persist")
	}
}

func TestStorageUpdateNonexistantPost(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	_, err := store.InsertPost(ctx, storage.InPost{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = store.UpdatePost(ctx, uuid.New(), storage.InPost{})
	if !errors.Is(err, storage.ErrPostNotFound) {
		t.Fatalf("error: %v", err)
	}
}

func TestStorageUpdateDeletedPost(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	post, err := store.InsertPost(ctx, storage.InPost{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = store.DeletePost(ctx, post.Id)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = store.UpdatePost(ctx, post.Id, storage.InPost{})
	if !errors.Is(err, storage.ErrPostIsDeleted) {
		t.Fatalf("managed to update deleted post: %v", err)
	}
}

func TestStorageInsertComment(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	inPost := storage.InPost{
		UserId:  uuid.New(),
		IsMute:  false,
		Content: "content",
	}

	post, err := store.InsertPost(ctx, inPost)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	inComm := storage.InComment{
		UserId:  uuid.New(),
		Content: "content",
	}

	comm, err := store.InsertComment(ctx, post.Id, nil, inComm)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if comm.InComment != inComm {
		t.Fatalf("comment has wrong data")
	}

	prevPost, err := store.GetPost(ctx, post.Id)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if _, ok := prevPost.Comments[comm.Id]; !ok {
		t.Fatalf("comment didn't persist")
	}
}

func TestStorageInsertCommentIntoNonexistantPost(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	_, err := store.InsertPost(ctx, storage.InPost{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = store.InsertComment(ctx, uuid.New(), nil, storage.InComment{})
	if !errors.Is(err, storage.ErrPostNotFound) {
		t.Fatalf("error: %v", err)
	}
}

func TestStorageInsertReply(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	post, err := store.InsertPost(ctx, storage.InPost{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	comm, err := store.InsertComment(ctx, post.Id, nil, storage.InComment{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	inRepl := storage.InComment{
		UserId:  uuid.New(),
		Content: "content",
	}

	repl, err := store.InsertComment(ctx, post.Id, &comm.Id, inRepl)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if repl.InComment != inRepl {
		t.Fatalf("reply has wrong data")
	}
}

func TestStorageInsertReplyIntoNonexistantComment(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	post, err := store.InsertPost(ctx, storage.InPost{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = store.InsertComment(ctx, post.Id, nil, storage.InComment{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	rnd := uuid.New()

	_, err = store.InsertComment(ctx, post.Id, &rnd, storage.InComment{})
	if !(errors.Is(err, storage.ErrCommNotFound)) {
		t.Fatalf("error: %v", err)
	}
}

func TestStorageDeleteComment(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	post, err := store.InsertPost(ctx, storage.InPost{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	comm, err := store.InsertComment(ctx, post.Id, nil, storage.InComment{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = store.DeleteComment(ctx, post.Id, comm.Id)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestStorageDeleteReply(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	post, err := store.InsertPost(ctx, storage.InPost{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	comm, err := store.InsertComment(ctx, post.Id, nil, storage.InComment{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	repl, err := store.InsertComment(ctx, post.Id, &comm.Id, storage.InComment{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = store.DeleteComment(ctx, post.Id, repl.Id)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestStorageDeleteCommentInNonexistantPost(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	post, err := store.InsertPost(ctx, storage.InPost{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	comm, err := store.InsertComment(ctx, post.Id, nil, storage.InComment{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = store.DeleteComment(ctx, uuid.New(), comm.Id)
	if !errors.Is(err, storage.ErrPostNotFound) {
		t.Fatalf("error: %v", err)
	}
}

func TestStorageUpdateComment(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	post, err := store.InsertPost(ctx, storage.InPost{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	id := uuid.New()

	inComm := storage.InComment{
		UserId:  id,
		Content: "content",
	}

	comm, err := store.InsertComment(ctx, post.Id, nil, inComm)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	inUpd := storage.InComment{
		UserId:  id,
		Content: "123123123",
	}

	upd, err := store.UpdateComment(ctx, post.Id, comm.Id, inUpd)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if upd.InComment != inUpd {
		t.Fatalf("updates didn't persist")
	}
}

func TestStorageUpdateCommentInNonexistantPost(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	post, err := store.InsertPost(ctx, storage.InPost{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	comm, err := store.InsertComment(ctx, post.Id, nil, storage.InComment{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = store.UpdateComment(ctx, uuid.New(), comm.Id, storage.InComment{})
	if !errors.Is(err, storage.ErrPostNotFound) {
		t.Fatalf("error: %v", err)
	}
}

func TestStorageUpdateReply(t *testing.T) {
	ctx := context.Background()

	store, _ = mem.NewStorage(conf, nil, nil, nil)

	post, err := store.InsertPost(ctx, storage.InPost{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	comm, err := store.InsertComment(ctx, post.Id, nil, storage.InComment{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	id := uuid.New()

	inRepl := storage.InComment{
		UserId:  id,
		Content: "content",
	}

	repl, err := store.InsertComment(ctx, post.Id, &comm.Id, inRepl)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	inUpd := storage.InComment{
		UserId:  id,
		Content: "skibidi",
	}

	upd, err := store.UpdateComment(ctx, post.Id, repl.Id, inUpd)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if upd.InComment != inUpd {
		t.Fatalf("updates didn't persist")
	}
}
