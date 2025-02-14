package storage

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrSessionNotFound   = errors.New("session not found")
	ErrNotImplemented    = errors.New("not implemented")
	ErrRoleNotFound      = errors.New("role not found")
)
