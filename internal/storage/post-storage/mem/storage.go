package mem

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/cutlery47/posts/config"
	storage "github.com/cutlery47/posts/internal/storage/post-storage"
	"github.com/google/uuid"
)

type memStorage struct {
	mu *sync.RWMutex
	// PostId -> Post
	posts map[uuid.UUID]storage.Post

	// restore src / dump dst
	rfd, wfd *os.File

	conf config.PostStorage
}

func NewStorage(conf config.PostStorage, rfd, wfd *os.File, errChan chan<- error) (*memStorage, error) {
	var (
		ms = &memStorage{
			mu:    &sync.RWMutex{},
			posts: make(map[uuid.UUID]storage.Post),
			conf:  conf,
		}
	)

	if !conf.DumpEnabled {
		return ms, nil
	}

	ms.wfd = wfd
	ms.rfd = rfd

	if err := ms.restore(); err != nil {
		return nil, err
	}

	go ms.dump(errChan)

	return ms, nil
}

func (ms *memStorage) GetPost(ctx context.Context, id uuid.UUID) (*storage.Post, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.RLock()
	defer ms.mu.RUnlock()

	v, ok := ms.posts[id]
	if !ok {
		return nil, storage.ErrPostNotFound
	}

	return &v, nil
}

func (ms *memStorage) GetPosts(ctx context.Context) ([]storage.Post, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.RLock()
	defer ms.mu.RUnlock()

	posts := make([]storage.Post, 0, len(ms.posts))

	for _, v := range ms.posts {
		posts = append(posts, v)
	}

	return posts, nil
}

func (ms *memStorage) InsertPost(ctx context.Context, in storage.InPost) (*storage.Post, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	post := toPost(in)

	// loop until no collisions detected
	for _, ok := ms.posts[post.Id]; ok; _, ok = ms.posts[post.Id] {
		post = toPost(in)
	}

	ms.posts[post.Id] = post

	return &post, nil
}

func (ms *memStorage) DeletePost(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	post, ok := ms.posts[id]
	if !ok {
		return nil, storage.ErrPostNotFound
	}

	if post.DeletedAt != nil {
		return nil, storage.ErrPostIsDeleted
	}

	ts := time.Now()

	post.DeletedAt = &ts
	ms.posts[id] = post

	return &id, nil
}

func (ms *memStorage) UpdatePost(ctx context.Context, id uuid.UUID, in storage.InPost) (*storage.Post, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	post, ok := ms.posts[id]
	if !ok {
		return nil, storage.ErrPostNotFound
	}

	if post.DeletedAt != nil {
		return nil, storage.ErrPostIsDeleted
	}

	post.UpdatedAt = time.Now()
	post.Content = in.Content
	post.IsMute = in.IsMute

	ms.posts[id] = post

	return &post, nil
}

func (ms *memStorage) GetComment(ctx context.Context, postId, commentId uuid.UUID) (*storage.Comment, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.RLock()
	defer ms.mu.RUnlock()

	post, ok := ms.posts[postId]
	if !ok {
		return nil, storage.ErrPostNotFound
	}

	comm, ok := getComment(post, commentId)
	if !ok {
		return nil, storage.ErrCommNotFound
	}

	return comm, nil
}

func (ms *memStorage) InsertComment(ctx context.Context, postId uuid.UUID, parentId *uuid.UUID, in storage.InComment) (*storage.Comment, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	post, ok := ms.posts[postId]
	if !ok {
		return nil, storage.ErrPostNotFound
	}

	// fast path
	// entered if a given comment is on the first level of tree depth (O(1): single insertion)
	if parentId == nil {
		comm, err := insertComment(post, in)
		if err != nil {
			return nil, err
		}
		return comm, nil
	}

	// slow path
	// entered if a given comment is somwhere down the tree (O(N): full dfs traversal + insertion)
	parent, ok := getComment(post, *parentId)
	if !ok {
		return nil, storage.ErrCommNotFound
	}

	return insertReply(*parent, in)
}

func (ms *memStorage) UpdateComment(ctx context.Context, postId, commentId uuid.UUID, in storage.InComment) (*storage.Comment, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	post, ok := ms.posts[postId]
	if !ok {
		return nil, storage.ErrPostNotFound
	}

	return updateComment(post, commentId, in)
}

func (ms *memStorage) DeleteComment(ctx context.Context, postId, commentId uuid.UUID) (*uuid.UUID, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	post, ok := ms.posts[postId]
	if !ok {
		return nil, storage.ErrPostNotFound
	}

	return deleteComment(post, commentId)
}

// dump current state of the storage into given io.ReadWriter
func (ms *memStorage) dump(errChan chan<- error) {
	for {
		time.Sleep(ms.conf.DumpInterval)

		err := func() error {
			ms.mu.Lock()
			defer ms.mu.Unlock()

			// clear prev contents
			if err := ms.wfd.Truncate(0); err != nil {
				return err
			}

			// move pointer to the beginning
			if _, err := ms.wfd.Seek(0, 0); err != nil {
				return err
			}

			// flush storage state
			err := json.NewEncoder(ms.wfd).Encode(ms.posts)
			if err != nil {
				return err
			}

			return nil
		}()

		if err != nil {
			errChan <- fmt.Errorf("%v: %v", ErrBadDump, err)
			break
		}
	}

}

// restores last state of the storage from given io.ReadWriter
func (ms *memStorage) restore() error {
	err := json.NewDecoder(ms.rfd).Decode(&ms.posts)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("%v: %v", ErrBadRestore, err)
	}

	return nil
}
