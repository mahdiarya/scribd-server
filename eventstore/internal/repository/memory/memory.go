package memory

import (
	"context"
	"sync"

	"scribd/eventstore/internal/repository"
	"scribd/eventstore/pkg/model"
)

// Repository defines a memory movie matadata repository.
type Repository struct {
	sync.RWMutex
	data map[string]*model.Event
}

// New creates a new memory repository.
func New() *Repository {
	return &Repository{data: map[string]*model.Event{}}
}

// Get retrieves movie metadata for by movie id.
func (r *Repository) Get(_ context.Context, id string) (*model.Event, error) {
	r.RLock()
	defer r.RUnlock()
	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(_ context.Context, id string, metadata *model.Event) error {
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return nil
}
