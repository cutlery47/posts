package postgres

import storage "github.com/cutlery47/posts/internal/storage/user-storage"

func NewStorage() (storage.Storage, error) {
	return nil, storage.ErrNotImplemented
}
