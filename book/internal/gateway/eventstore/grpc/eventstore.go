package grpc

import (
	"context"
	"errors"

	"scribd/eventstore/pkg/model"
	gen "scribd/gen/proto/v1"
	"scribd/internal/grpcutil"
	"scribd/pkg/discovery"
)

// Gateway defines an gRPC gateway for a rating service.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for a rating service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// GetAggregatedRating returns the aggregated rating for a record or ErrNotFound if there are no ratings for it.
func (g *Gateway) GetEvents(ctx context.Context, eventID string, aggregateID string) ([]*gen.Event, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := gen.NewEventStoreClient(conn)
	resp, err := client.GetEvents(ctx, &gen.GetEventsRequest{EventId: eventID, AggregateId: aggregateID})
	if err != nil {
		return nil, err
	}
	return resp.Events, nil
}

func (g *Gateway) CreateEvent(ctx context.Context, event model.Event) (bool, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	client := gen.NewEventStoreClient(conn)
	resp, err := client.CreateEvent(ctx, &gen.CreateEventRequest{Event: &gen.Event{
		EventId:       event.EventID,
		EventType:     event.EventType,
		AggregateId:   event.AggregateID,
		AggregateType: event.AggregateType,
		EventData:     event.EventData,
		Stream:        event.Stream,
	}})
	if err != nil {
		return false, err
	}
	return resp.IsSuccess, errors.New(resp.Error)
}
