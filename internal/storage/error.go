package storage

import "errors"

var (
	ErrPostNotFound   = errors.New("post not found")
	ErrPostIsDeleted  = errors.New("post has been deleted")
	ErrPostIsMute     = errors.New("post is mute")
	ErrCommNotFound   = errors.New("comment not found")
	ErrCommIsDeleted  = errors.New("comment has been deleted")
	ErrNotImplemented = errors.New("not implemented")
)
