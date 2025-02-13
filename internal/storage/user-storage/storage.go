package storage

import "context"

type Storage interface {
	// registeres user with given inputs
	Register(ctx context.Context, in InUser) (*User, error)
	// logs given user in
	Login(ctx context.Context, user User) (*Session, error)
	// logs given user out
	Logout(ctx context.Context, sesh Session) error
}
