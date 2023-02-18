package grpc

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"scribd/eventstore/internal/repository/memory"
	"scribd/eventstore/pkg/model"
	"scribd/eventstore/pkg/natsutil"
	gen "scribd/gen/proto/v1"
)

type Server struct {
	gen.UnimplementedEventStoreServer
	Repository *memory.Repository
	Nats       *natsutil.NATSComponent
}

// publishEvent publishes an event via NATS JetStream server
func publishEvent(component *natsutil.NATSComponent, event *gen.Event) {
	// Creates JetStreamContext to publish messages into JetStream Stream
	jetStreamContext, _ := component.JetStreamContext()
	subject := event.EventType
	eventMsg := []byte(event.EventData)
	// Publish message on subject (channel)
	jetStreamContext.Publish(subject, eventMsg)
	log.Println("Published message on subject: " + subject)
}

// CreateEvent creates a new event into the event store
func (s *Server) CreateEvent(ctx context.Context, eventRequest *gen.CreateEventRequest) (*gen.CreateEventResponse, error) {
	err := s.Repository.Put(ctx, uuid.New().String(), &model.Event{
		EventID:       eventRequest.Event.EventId,
		EventType:     eventRequest.Event.EventType,
		AggregateID:   eventRequest.Event.AggregateId,
		AggregateType: eventRequest.Event.AggregateType,
		EventData:     eventRequest.Event.EventData,
		Stream:        eventRequest.Event.Stream,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	log.Println("Event is created")
	go publishEvent(s.Nats, eventRequest.Event)
	return &gen.CreateEventResponse{IsSuccess: true, Error: ""}, nil
}

// GetEvents gets all events for the given aggregate and event
func (s *Server) GetEvents(ctx context.Context, filter *gen.GetEventsRequest) (*gen.GetEventsResponse, error) {
	events, err := s.Repository.Get(ctx, filter.EventId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &gen.GetEventsResponse{Events: []*gen.Event{{
		EventId:       events.EventID,
		EventType:     events.EventType,
		AggregateId:   events.AggregateID,
		AggregateType: events.AggregateType,
		EventData:     events.EventData,
		Stream:        events.Stream,
	}}}, nil
}

// GetEventsStream get stream of events for the given event
func (s *Server) GetEventsStream(*gen.GetEventsRequest, gen.EventStore_GetEventsStreamServer) error {
	return nil
}
