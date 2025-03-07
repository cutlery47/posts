package storage

import (
	"context"

	"github.com/google/uuid"
)

type Storage interface {
	// retrieves a single post by provided id
	GetPost(ctx context.Context, id uuid.UUID) (*Post, error)
	// retrieves all posts
	GetPosts(ctx context.Context) ([]Post, error)
	// inserts a single post
	InsertPost(ctx context.Context, in InPost) (*Post, error)
	// deletes a single post by provided id
	DeletePost(ctx context.Context, id uuid.UUID) (*uuid.UUID, error)
	// updates a single post by provided id
	UpdatePost(ctx context.Context, id uuid.UUID, in InPost) (*Post, error)

	GetComment(ctx context.Context, postId, commentId uuid.UUID) (*Comment, error)
	// inserts a single comment for a post by provided id
	InsertComment(ctx context.Context, postId uuid.UUID, parentId *uuid.UUID, in InComment) (*Comment, error)
	// deletes a single comment for a post by provided id
	DeleteComment(ctx context.Context, postId, commentId uuid.UUID) (*uuid.UUID, error)
	// updates a single comment for a post by provided id
	UpdateComment(ctx context.Context, postId, commentId uuid.UUID, in InComment) (*Comment, error)
}
