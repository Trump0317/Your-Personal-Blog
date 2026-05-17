package db

import (
	"context"
	"errors"
	"sync"

	"github.com/ypb/your-personal-blog/internal/hub/model"
)

type MemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

func NewMemoryUserRepo() *MemoryUserRepo {
	return &MemoryUserRepo{
		users: make(map[string]*model.User),
	}
}

func (r *MemoryUserRepo) Create(ctx context.Context, user model.User) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.APIKey] = &user
	return user.APIKey, nil
}

func (r *MemoryUserRepo) GetByAPIKey(ctx context.Context, apiKey string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[apiKey]
	if !ok {
		return nil, errors.New("user not found by APIKey")
	}
	return user, nil
}

func (r *MemoryUserRepo) Delete(ctx context.Context, apiKey string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.users, apiKey)
	return nil
}

func (r *MemoryUserRepo) UpdateStatus(ctx context.Context, apiKey string, status model.UserStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[apiKey]
	if !ok {
		return errors.New("user not found")
	}
	user.Status = status
	return nil
}
