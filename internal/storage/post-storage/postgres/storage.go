package postgres

import storage "github.com/cutlery47/posts/internal/storage/post-storage"

func NewStorage() (storage.Storage, error) {
	return nil, storage.ErrNotImplemented
}
