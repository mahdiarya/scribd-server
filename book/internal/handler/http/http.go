package http

import (
	"errors"
	"log"
	"net/http"

	"scribd/book/internal/features"
	"scribd/book/internal/repository"
	"scribd/book/pkg/model"

	"github.com/gin-gonic/gin"
)

// Handler defines a book metadata HTTP handler.
type Handler struct {
	ctrl *features.Application
}

// New creates a new book metadata HTTP handler.
func New(ctrl *features.Application) *Handler {
	return &Handler{ctrl}
}

// HandleRoutes handles GET /book requests.
func (h *Handler) HandleRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	// Command Routes
	v1.POST("/book", h.PostBook)
	// Queries Routes
	v1.GET("/book", h.GetBook)
}

// PostBook handles POST /book requests.
func (h *Handler) PostBook(context *gin.Context) {
	body := model.Book{}
	// using BindJson method to serialize body with struct
	if err := context.BindJSON(&body); err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := h.ctrl.CreateBook(context.Request.Context(), &body)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		context.JSON(http.StatusInternalServerError, &body)
		return
	}
	context.JSON(http.StatusAccepted, &res)
}

// GetMovieDetails handles GET /book requests.
func (h *Handler) GetBook(context *gin.Context) {
	id := context.Query("id")
	if id == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "id invalid",
		})
		return
	}
	m, err := h.ctrl.GetBook(context.Request.Context(), id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		context.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "book not found",
		})
		return
	} else if err != nil {
		log.Printf("Repository get error: %v\n", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "internal error",
		})
		return
	}
	context.JSON(http.StatusAccepted, gin.H{
		"code":    http.StatusAccepted,
		"message": "internal error",
		"data":    m,
	})
}
