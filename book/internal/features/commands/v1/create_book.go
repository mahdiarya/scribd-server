package v1

import (
	"context"

	"scribd/book/pkg/model"
)

type bookRepository interface {
	Get(ctx context.Context, id string) (*model.Book, error)
	Put(ctx context.Context, id string, m *model.Book) error
}

// Controller defines a book service controller.
type CreateBookController struct {
	repo bookRepository
}

// New creates a book service controller.
func NewCreateBookController(repo bookRepository) CreateBookController {
	return CreateBookController{repo}
}

func (h CreateBookController) CreateOrder(ctx context.Context, cmd *model.Book) error {
	return h.repo.Put(ctx, cmd.ID, cmd)
}
