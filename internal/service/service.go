package service

import (
	"context"
	"errors"
	"slices"

	post "github.com/cutlery47/posts/internal/storage/post-storage"
	storage "github.com/cutlery47/posts/internal/storage/post-storage"
	user "github.com/cutlery47/posts/internal/storage/user-storage"
	"github.com/google/uuid"
)

var (
	SortNewest    = "newest"
	SortOldest    = "oldest"
	SortUpvotes   = "upvoted"
	SortDownvotes = "downvoted"
)

type Service struct {
	ps post.Storage
	us user.Storage
}

func New(ps post.Storage, us user.Storage) (*Service, error) {
	return &Service{
		ps: ps,
		us: us,
	}, nil
}

func (s *Service) GetPost(ctx context.Context, id uuid.UUID) (*post.Post, error) {
	return s.ps.GetPost(ctx, id)
}

func (s *Service) GetPosts(ctx context.Context, limit *int, offset *int, sortBy string) ([]storage.Post, error) {
	posts, err := s.ps.GetPosts(ctx)
	if err != nil {
		return nil, err
	}

	posts, err = s.sortPosts(posts, sortBy)
	if err != nil {
		return nil, err
	}

	if offset != nil {
		posts = posts[min(*offset, len(posts)):]
	}

	if limit != nil {
		posts = posts[:min(*limit, len(posts))]
	}

	return posts, nil
}

func (s *Service) InsertPost(ctx context.Context, in post.InPost) (*post.Post, error) {
	return s.ps.InsertPost(ctx, in)
}

func (s *Service) DeletePost(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	// pass userId somehow + validation

	return s.ps.DeletePost(ctx, id)
}

func (s *Service) UpdatePost(ctx context.Context, id uuid.UUID, in post.InPost) (*post.Post, error) {
	// validation

	return s.ps.UpdatePost(ctx, id, in)
}

func (s *Service) InsertComment(ctx context.Context, postId uuid.UUID, parentId *uuid.UUID, in post.InComment) (*post.Comment, error) {
	// vld

	return s.ps.InsertComment(ctx, postId, parentId, in)
}

func (s *Service) DeleteComment(ctx context.Context, postId, commentId uuid.UUID) (*uuid.UUID, error) {
	// vid

	return s.ps.DeleteComment(ctx, postId, commentId)
}

func (s *Service) sortPosts(posts []post.Post, sortBy string) ([]post.Post, error) {
	switch sortBy {
	case SortNewest:
		slices.SortFunc(posts, func(a, b post.Post) int {
			return b.CreatedAt.Compare(a.CreatedAt)
		})
	case SortOldest:
		slices.SortFunc(posts, func(a, b post.Post) int {
			return a.CreatedAt.Compare(b.CreatedAt)
		})
	case SortUpvotes:
		slices.SortFunc(posts, func(a, b post.Post) int {
			return int(b.Upvotes) - int(a.Upvotes)
		})
	case SortDownvotes:
		slices.SortFunc(posts, func(a, b post.Post) int {
			return int(a.Upvotes) - int(b.Upvotes)
		})
	default:
		return nil, errors.New("undefined sort key")
	}

	return posts, nil
}
