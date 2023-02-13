package v1

import (
	"context"
	"errors"

	"scribd/book/internal/repository"
	"scribd/book/pkg/model"
)

// ErrNotFound is returned when a requested record is not found.
var ErrNotFound = errors.New("not found")

type bookRepository interface {
	Get(ctx context.Context, id string) (*model.Book, error)
	Put(ctx context.Context, id string, m *model.Book) error
}

// Controller defines a book service controller.
type GetBookController struct {
	repo bookRepository
}

// New creates a book service controller.
func NewGetBookController(repo bookRepository) GetBookController {
	return GetBookController{repo}
}

func (h GetBookController) GetBook(ctx context.Context, id string) (*model.Book, error) {
	res, err := h.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}
