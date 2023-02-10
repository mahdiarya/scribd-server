package http

import (
	"net/http"
	"scribd/book/internal/controller/book"
)

// Handler defines a book metadata HTTP handler.
type Handler struct {
	ctrl *book.Controller
}

// New creates a new book metadata HTTP handler.
func New(ctrl *book.Controller) *Handler {
	return &Handler{ctrl}
}

// GetMovieDetails handles GET /book requests.
func (h *Handler) GetBookDetails(w http.ResponseWriter, req *http.Request) {

}
