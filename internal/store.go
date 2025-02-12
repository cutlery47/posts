package internal

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Storage interface {
	// retrieves a single post by provided id
	GetPost(ctx context.Context, id uuid.UUID) (*Post, error)
	// retrieves all posts
	GetPosts(ctx context.Context) []Post
	// inserts a single post
	InsertPost(ctx context.Context, in InPost) (*Post, error)
	// deletes a single post by provided id
	DeletePost(ctx context.Context, id uuid.UUID) error
	// updates a single post by provided id
	UpdatePost(ctx context.Context, id uuid.UUID, in InPost) (*Post, error)

	// inserts a single comment for a post by provided id
	InsertComment(ctx context.Context, postId uuid.UUID, parentId *uuid.UUID, in InComment) (*Comment, error)
	// deletes a single comment for a post by provided id
	DeleteComment(ctx context.Context, postId, commentId uuid.UUID) error
	// updates a single comment for a post by provided id
	UpdateComment(ctx context.Context, postId, commentId uuid.UUID, in InComment) (*Comment, error)
}

type memStorage struct {
	mu *sync.RWMutex

	// PostId -> Post
	posts map[uuid.UUID]Post
}

func NewMemStorage() *memStorage {
	return &memStorage{
		mu:    &sync.RWMutex{},
		posts: make(map[uuid.UUID]Post),
	}
}

func (ms *memStorage) GetPost(ctx context.Context, id uuid.UUID) (*Post, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.RLock()
	defer ms.mu.RUnlock()

	v, ok := ms.posts[id]
	if !ok {
		return nil, ErrPostNotFound
	}

	return &v, nil
}

func (ms *memStorage) GetPosts(ctx context.Context) ([]Post, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.RLock()
	defer ms.mu.RUnlock()

	posts := make([]Post, 0, len(ms.posts))

	for _, v := range ms.posts {
		posts = append(posts, v)
	}

	return posts, nil
}

func (ms *memStorage) InsertPost(ctx context.Context, in InPost) (*Post, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	post := toPost(in)

	// loop until no collisions detected
	for _, ok := ms.posts[post.id]; ok; _, ok = ms.posts[post.id] {
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
		return ErrPostNotFound
	}

	ts := time.Now()

	post.deletedAt = &ts
	ms.posts[id] = post

	return nil
}

func (ms *memStorage) UpdatePost(ctx context.Context, id uuid.UUID, in InPost) (*Post, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	post, ok := ms.posts[id]
	if !ok {
		return nil, ErrPostNotFound
	}

	post.updatedAt = time.Now()
	post.content = in.content
	post.isMute = in.isMute

	ms.posts[id] = post

	return &post, nil
}

func (ms *memStorage) InsertComment(ctx context.Context, postId uuid.UUID, parentId *uuid.UUID, in InComment) (*Comment, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	post, ok := ms.posts[postId]
	if !ok {
		return nil, ErrPostNotFound
	}

	var (
		comm Comment
	)

	// fast path
	// entered if a given comment is on the first level of tree depth (O(1): single insertion)
	if parentId == nil {
		comm = post.insertComment(in)
		return &comm, nil
	}

	// slow path
	// entered if a given comment is somwhere down the tree (O(N): full dfs traversal + insertion)
	parent, ok := post.getComment(*parentId)
	if !ok {
		return nil, ErrCommNotFound
	}

	comm = parent.insertReply(in)

	return &comm, nil
}

func (ms *memStorage) UpdateComment(ctx context.Context, postId, commentId uuid.UUID, in InComment) (*Comment, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	post, ok := ms.posts[postId]
	if !ok {
		return nil, ErrPostNotFound
	}

	return post.updateComment(commentId, in)
}

func (ms *memStorage) DeleteComment(ctx context.Context, postId, commentId uuid.UUID) error {
	if err := ctxDone(ctx); err != nil {
		return err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	post, ok := ms.posts[postId]
	if !ok {
		return ErrPostNotFound
	}

	return post.deleteComment(commentId)
}
