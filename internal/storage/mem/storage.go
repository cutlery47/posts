package mem

import (
	"context"
	"sync"
	"time"

	"github.com/cutlery47/posts/internal/storage"
	"github.com/google/uuid"
)

type memStorage struct {
	mu *sync.RWMutex

	// PostId -> Post
	posts map[uuid.UUID]storage.Post
}

func NewMemStorage() *memStorage {
	return &memStorage{
		mu:    &sync.RWMutex{},
		posts: make(map[uuid.UUID]storage.Post),
	}
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

	return &post, nil
}

func (ms *memStorage) DeletePost(ctx context.Context, id uuid.UUID) error {
	if err := ctxDone(ctx); err != nil {
		return err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	post, ok := ms.posts[id]
	if !ok {
		return storage.ErrPostNotFound
	}

	if post.DeletedAt != nil {
		return storage.ErrPostIsDeleted
	}

	ts := time.Now()

	post.DeletedAt = &ts
	ms.posts[id] = post

	return nil
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

func (ms *memStorage) DeleteComment(ctx context.Context, postId, commentId uuid.UUID) error {
	if err := ctxDone(ctx); err != nil {
		return err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	post, ok := ms.posts[postId]
	if !ok {
		return storage.ErrPostNotFound
	}

	return deleteComment(post, commentId)
}
