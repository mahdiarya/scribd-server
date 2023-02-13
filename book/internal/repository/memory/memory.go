package memory

import (
	"context"
	"sync"

	"scribd/book/internal/repository"
	"scribd/book/pkg/model"
)

// Repository defines a memory movie matadata repository.
type Repository struct {
	sync.RWMutex
	data map[string]*model.Book
}

// New creates a new memory repository.
func New() *Repository {
	return &Repository{data: map[string]*model.Book{}}
}

// Get retrieves movie metadata for by movie id.
func (r *Repository) Get(_ context.Context, id string) (*model.Book, error) {
	r.RLock()
	defer r.RUnlock()
	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(_ context.Context, id string, metadata *model.Book) error {
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return nil
}
