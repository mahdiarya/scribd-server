package grpc

import (
	"context"
	"errors"
	"log"

	"scribd/book/internal/features"
	"scribd/book/internal/repository"
	"scribd/book/pkg/model"
	gen "scribd/gen/proto/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a movie gRPC handler.
type Handler struct {
	gen.UnimplementedBookServiceServer
	ctrl *features.Application
}

// New creates a new movie gRPC handler.
func New(ctrl *features.Application) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) PostBook(ctx context.Context, req *gen.PostBookRequest) (*gen.PostBookResponse, error) {
	if req == nil || req.Title == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty title")
	}
	res, err := h.ctrl.CreateBook(ctx, &model.Book{Title: req.Title})
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.PostBookResponse{
		Book: &gen.Book{
			Id:    res.ID,
			Title: res.Title,
		}}, nil
}

// GetMovieDetails returns moviie details by id.
func (h *Handler) GetBook(ctx context.Context, req *gen.GetBookRequest) (*gen.GetBookResponse, error) {
	if req == nil || req.BookId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	m, err := h.ctrl.GetBook(ctx, req.BookId)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.GetBookResponse{
		Book: &gen.Book{
			Id:    m.ID,
			Title: m.Title,
		},
	}, nil
}
