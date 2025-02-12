package internal

import "errors"

var (
	ErrPostNotFound   = errors.New("post not found")
	ErrCommNotFound   = errors.New("comment not found")
	ErrNotImplemented = errors.New("not implemented")
)
