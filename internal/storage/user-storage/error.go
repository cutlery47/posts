package storage

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrSessionNotFound   = errors.New("session not found")
	ErrNotImplemented    = errors.New("not implemented")
)
