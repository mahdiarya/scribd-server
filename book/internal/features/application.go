package features

import (
	"context"

	commandsV1 "scribd/book/internal/features/commands/v1"
	queriesV1 "scribd/book/internal/features/queries/v1"
	"scribd/book/pkg/model"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		CreateBook(ctx context.Context, cmd model.Book) error
	}
	Queries interface {
		GetBook(ctx context.Context, id string) (*model.Book, error)
	}
	bookRepository interface {
		Get(ctx context.Context, id string) (*model.Book, error)
		Put(ctx context.Context, id string, m *model.Book) error
	}

	Application struct {
		appCommands
		appQueries
	}
	appCommands struct {
		commandsV1.CreateBookController
	}
	appQueries struct {
		queriesV1.GetBookController
	}
)

func New(repo bookRepository) *Application {
	return &Application{
		appCommands: appCommands{
			CreateBookController: commandsV1.NewCreateBookController(repo),
		},
		appQueries: appQueries{
			GetBookController: queriesV1.NewGetBookController(repo),
		},
	}
}
