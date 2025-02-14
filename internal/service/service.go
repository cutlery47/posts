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

func (s *Service) ListPosts(ctx context.Context, limit *int, offset *int, sortBy string) ([]storage.Post, error) {
	posts, err := s.ps.GetPosts(ctx)
	if err != nil {
		return nil, err
	}

	posts, err = s.sortPosts(posts, sortBy)
	if err != nil {
		return nil, err
	}

	if offset != nil {
		posts = posts[*offset:]
	}

	if limit != nil {
		posts = posts[:*limit]
	}

	return posts, nil
}

func (s *Service) sortPosts(posts []post.Post, sortBy string) ([]post.Post, error) {
	switch sortBy {
	case SortNewest:
		slices.SortFunc(posts, func(a, b post.Post) int {
			return a.CreatedAt.Compare(b.CreatedAt)
		})
	case SortOldest:
		slices.SortFunc(posts, func(a, b post.Post) int {
			return b.CreatedAt.Compare(a.CreatedAt)
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
