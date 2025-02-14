package service

import "errors"

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrWrongUserId    = errors.New("you pretending to be another user")
	ErrAccessDenied   = errors.New("access denied")
)
