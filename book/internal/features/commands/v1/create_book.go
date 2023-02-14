package v1

import (
	"context"

	"scribd/book/pkg/model"

	"github.com/google/uuid"
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

func (h CreateBookController) CreateBook(ctx context.Context, cmd *model.Book) (*model.Book, error) {
	cmd.ID = uuid.New().String()
	err := h.repo.Put(ctx, cmd.ID, cmd)
	if err != nil {
		return nil, err
	}
	return cmd, err
}
