package book

import (
	"context"
	"scribd/book/pkg/model"
)

// type commentGateway interface {
// 	Get(ctx context.Context, id string) (*metadatamodel.Metadata, error)
// }

type bookRepository interface {
	Get(ctx context.Context, id string) (*model.Book, error)
	Put(ctx context.Context, id string, m *model.Book) error
}

// Controller defines a book service controller.
type Controller struct {
	repo bookRepository
}

// New creates a book service controller.
func New(repo bookRepository) *Controller {
	return &Controller{repo}
}

// Get returns book metadata by id.
func (c *Controller) Get(ctx context.Context, id string) (*model.Book, error) {

	return nil, nil
}

// Put writes book metadata to repository.
func (c *Controller) Put(ctx context.Context, m *model.Book) error {
	return nil
}
