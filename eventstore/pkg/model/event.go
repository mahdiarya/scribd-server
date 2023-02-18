package model

type Event struct {
	EventID       string
	EventType     string
	AggregateID   string
	AggregateType string
	EventData     string
	Stream        string
}
