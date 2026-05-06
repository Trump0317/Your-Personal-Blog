package db

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ypb/your-personal-blog/internal/hub/model"
)

// MemoryFileRepo 是 FileRepo 接口的内存简单实现
type MemoryFileRepo struct {
	mu     sync.RWMutex
	files  map[string]*model.File
	nextID int64
}

func NewMemoryFileRepo() *MemoryFileRepo {
	return &MemoryFileRepo{
		files:  make(map[string]*model.File),
		nextID: 1,
	}
}

func (r *MemoryFileRepo) Create(ctx context.Context, file model.File) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := fmt.Sprintf("%d", r.nextID)
	r.nextID++

	newFile := file
	newFile.ID = id
	r.files[id] = &newFile

	return id, nil
}

func (r *MemoryFileRepo) GetByID(ctx context.Context, id string) (*model.File, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	file, ok := r.files[id]
	if !ok {
		return nil, errors.New("file not found in memory")
	}
	return file, nil
}

func (r *MemoryFileRepo) ListByUser(ctx context.Context, userID string) ([]*model.File, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.File
	for _, f := range r.files {
		if f.UploaderID == userID {
			result = append(result, f)
		}
	}
	return result, nil
}

func (r *MemoryFileRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.files, id)
	return nil
}

func (r *MemoryFileRepo) UpdateStatus(ctx context.Context, id string, status model.FileStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, ok := r.files[id]
	if !ok {
		return errors.New("file not found in memory")
	}
	file.Status = status
	return nil
}
