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

func (s *Service) GetSessionUser(ctx context.Context, seshId uuid.UUID) (uuid.UUID, error) {
	sesh, err := s.us.GetSession(ctx, seshId)
	if err != nil && errors.Is(err, user.ErrSessionNotFound) {
		return uuid.UUID{}, err
	}
	return sesh.UserId, nil
}

func (s *Service) Register(ctx context.Context, in user.InUser) (*user.User, error) {
	return s.us.Register(ctx, in)
}

func (s *Service) Login(ctx context.Context, in user.InUser) (*user.Session, error) {
	return s.us.Login(ctx, in)
}

func (s *Service) Logout(ctx context.Context, sesh user.Session) error {
	return s.us.Logout(ctx, sesh)
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

func (s *Service) InsertPost(ctx context.Context, in post.InPost, userId uuid.UUID) (*post.Post, error) {
	if in.UserId != userId {
		return nil, ErrWrongUserId
	}

	return s.ps.InsertPost(ctx, in)
}

func (s *Service) DeletePost(ctx context.Context, id uuid.UUID, userId uuid.UUID) (*uuid.UUID, error) {
	post, err := s.ps.GetPost(ctx, id)
	if err != nil {
		return nil, err
	}

	if post.UserId != userId {
		return nil, ErrAccessDenied
	}

	return s.ps.DeletePost(ctx, id)
}

func (s *Service) UpdatePost(ctx context.Context, id, userId uuid.UUID, in post.InPost) (*post.Post, error) {
	post, err := s.ps.GetPost(ctx, id)
	if err != nil {
		return nil, err
	}

	if post.UserId != userId {
		return nil, ErrAccessDenied
	}

	if in.UserId != userId {
		return nil, ErrWrongUserId
	}

	return s.ps.UpdatePost(ctx, id, in)
}

func (s *Service) InsertComment(ctx context.Context, postId, userId uuid.UUID, parentId *uuid.UUID, in post.InComment) (*post.Comment, error) {
	if in.UserId != userId {
		return nil, ErrWrongUserId
	}

	return s.ps.InsertComment(ctx, postId, parentId, in)
}

func (s *Service) DeleteComment(ctx context.Context, postId, commentId, userId uuid.UUID) (*uuid.UUID, error) {
	comm, err := s.ps.GetComment(ctx, postId, commentId)
	if err != nil {
		return nil, err
	}

	if comm.UserId != userId {
		return nil, ErrAccessDenied
	}

	return s.ps.DeleteComment(ctx, postId, commentId)
}

func (s *Service) UpdateComment(ctx context.Context, postId, commentId, userId uuid.UUID, in post.InComment) (*post.Comment, error) {
	comm, err := s.ps.GetComment(ctx, postId, commentId)
	if err != nil {
		return nil, err
	}

	if userId != comm.UserId {
		return nil, ErrAccessDenied
	}

	if userId != in.UserId {
		return nil, ErrWrongUserId
	}

	return s.ps.UpdateComment(ctx, postId, commentId, in)
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
