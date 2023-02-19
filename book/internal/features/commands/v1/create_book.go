package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"scribd/book/pkg/model"
	eventmodel "scribd/eventstore/pkg/model"

	"github.com/google/uuid"
	"google.golang.org/grpc/status"
)

const (
	eventType     = "BOOKS.created"
	aggregateType = "book"
	grpcUri       = "localhost:50051"
)

type eventstoreGateway interface {
	CreateEvent(ctx context.Context, event eventmodel.Event) (bool, error)
}

type bookRepository interface {
	Get(ctx context.Context, id string) (*model.Book, error)
	Put(ctx context.Context, id string, m *model.Book) error
}

// Controller defines a book service controller.
type CreateBookController struct {
	repo         bookRepository
	eventGateway eventstoreGateway
}

// New creates a book service controller.
func NewCreateBookController(repo bookRepository, eventGateway eventstoreGateway) CreateBookController {
	return CreateBookController{repo, eventGateway}
}

func (c CreateBookController) CreateBook(ctx context.Context, cmd *model.Book) (*model.Book, error) {
	cmd.ID = uuid.New().String()
	err := c.repo.Put(ctx, cmd.ID, cmd)
	if err != nil {
		return nil, err
	}
	return cmd, err
}

func (c CreateBookController) PublishCreateBook(ctx context.Context, cmd *model.Book) error {
	aggregateid, _ := uuid.NewUUID()
	eventid, _ := uuid.NewUUID()
	cmd.ID = aggregateid.String()
	bookJSON, _ := json.Marshal(cmd)

	event := eventmodel.Event{
		EventID:       eventid.String(),
		EventType:     eventType,
		AggregateID:   cmd.ID,
		AggregateType: aggregateType,
		EventData:     string(bookJSON),
		Stream:        "BOOKS",
	}

	var wg sync.WaitGroup
	wg.Add(1)
	var isSuccess bool
	var getEventErr error
	go func() {
		defer wg.Done()
		isSuccess, getEventErr = c.eventGateway.CreateEvent(ctx, event)
	}()
	wg.Wait()
	if err := getEventErr; err != nil {
		if st, ok := status.FromError(err); ok {
			return fmt.Errorf("error from RPC server with: status code:%s message:%s", st.Code().String(), st.Message())
		}
		return fmt.Errorf("error from RPC server: %w", err)
	}
	if isSuccess {
		return nil
	}
	return errors.New("error from RPC server")
}
