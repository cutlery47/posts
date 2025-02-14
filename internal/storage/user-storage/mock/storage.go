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

func NewStorage() *mockStorage {
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

	if in.Role != storage.AdminRole && in.Role != storage.UserRole {
		return nil, storage.ErrRoleNotFound
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	for _, v := range ms.users {
		if v.Name == in.Name {
			return nil, storage.ErrUserAlreadyExists
		}
	}

	user := toUser(in)
	for _, ok := ms.users[user.Id]; ok; _, ok = ms.users[user.Id] {
		user = toUser(in)
	}

	ms.users[user.Id] = user

	return &user, nil
}

func (ms *mockStorage) Login(ctx context.Context, in storage.InUser) (*storage.Session, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	var (
		user *storage.User
	)

	// searching if user with provided name exists
	for _, v := range ms.users {
		if v.Name == in.Name {
			user = &v
			break
		}
	}

	if user == nil {
		return nil, storage.ErrUserNotFound
	}

	// imagine some password validation here, etc...

	return newSession(user.Id, time.Now().Add(ms.conf.SessionDuration)), nil
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

func (ms *mockStorage) GetSession(ctx context.Context, id uuid.UUID) (*storage.Session, error) {
	if err := ctxDone(ctx); err != nil {
		return nil, err
	}

	ms.mu.RLock()
	defer ms.mu.RUnlock()

	s, ok := ms.sessions[id]
	if !ok {
		return nil, storage.ErrSessionNotFound
	}

	return &s, nil
}
