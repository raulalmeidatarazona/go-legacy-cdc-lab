package cdc

import (
	"encoding/json"
	"errors"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// OutboxRow maps the append-only outbox table, not a customer's schema.
type OutboxRow struct {
	SequenceID uint64
	EventID    string
	EventType  string
	Aggregate  string
	Payload    json.RawMessage
}

func FromBinlogRow(values []any) (OutboxRow, error) {
	if len(values) != 6 {
		return OutboxRow{}, fmt.Errorf("outbox row has %d columns, want 6", len(values))
	}
	var sequence uint64
	switch v := values[0].(type) {
	case uint64:
		sequence = v
	case int64:
		if v < 0 {
			return OutboxRow{}, errors.New("sequence_id is negative")
		}
		sequence = uint64(v)
	default:
		return OutboxRow{}, fmt.Errorf("sequence_id has type %T", values[0])
	}
	id, err := asString(values[1])
	if err != nil {
		return OutboxRow{}, fmt.Errorf("event_id: %w", err)
	}
	kind, err := asString(values[2])
	if err != nil {
		return OutboxRow{}, fmt.Errorf("event_type: %w", err)
	}
	aggregate, err := asString(values[3])
	if err != nil {
		return OutboxRow{}, fmt.Errorf("aggregate_id: %w", err)
	}
	payload, err := asString(values[4])
	if err != nil {
		return OutboxRow{}, fmt.Errorf("payload: %w", err)
	}
	if !json.Valid([]byte(payload)) {
		return OutboxRow{}, errors.New("outbox payload is not valid JSON")
	}
	return OutboxRow{sequence, id, kind, aggregate, json.RawMessage(payload)}, nil
}

func asString(value any) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return "", fmt.Errorf("unexpected type %T", value)
	}
}

func (r OutboxRow) CloudEvent() cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetID(r.EventID)
	event.SetSource("/synthetic-legacy/inventory")
	event.SetType(r.EventType)
	event.SetSubject(r.Aggregate)
	event.SetDataContentType("application/json")
	event.SetExtension("originsequence", r.SequenceID)
	_ = event.SetData("application/json", r.Payload)
	return event
}
