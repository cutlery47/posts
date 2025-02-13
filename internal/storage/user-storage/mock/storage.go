package mock

import (
	"context"
	"sync"
	"time"

	"github.com/cutlery47/posts/config"
	storage "github.com/cutlery47/posts/internal/storage/user-storage"
	"github.com/google/uuid"
)

type mockStorage struct {
	users    map[uuid.UUID]storage.User
	sessions map[uuid.UUID]storage.Session

	mu *sync.RWMutex

	conf config.UserStorage
}

func NewMockStorage() *mockStorage {
	return &mockStorage{
		users:    make(map[uuid.UUID]storage.User),
		sessions: make(map[uuid.UUID]storage.Session),
		mu:       &sync.RWMutex{},
	}
}

func (ms *mockStorage) Register(ctx context.Context, in storage.InUser) (*storage.User, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	user := toUser(in)
	if _, ok := ms.users[user.Id]; ok {
		return nil, storage.ErrUserAlreadyExists
	}

	ms.users[user.Id] = user

	return &user, nil
}

func (ms *mockStorage) Login(ctx context.Context, user storage.User) (*storage.Session, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	// imagine some validation here...

	var (
		sesh *storage.Session
	)

	for _, v := range ms.sessions {
		if v.UserId == user.Id && v.ExpiresAt.Before(time.Now()) {
			sesh = &v
			break
		}
	}

	if sesh == nil {
		s := newSession(user.Id, time.Now().Add(ms.conf.SessionDuration))
		sesh = &s
	}

	return sesh, nil
}

func (ms *mockStorage) Logout(ctx context.Context, sesh storage.Session) error {
	if err := ctxDone(ctx); err != nil {
		return err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, ok := ms.sessions[sesh.Id]; !ok {
		return storage.ErrSessionNotFound
	}

	delete(ms.sessions, sesh.Id)

	return nil
}
