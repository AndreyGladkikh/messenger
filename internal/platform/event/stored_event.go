package event

import (
	"bytes"
	"messenger/messenger/internal/adapters/out/postgres/mapping"
	"messenger/messenger/internal/platform/db"
	"time"
)

type status string

const (
	statusPending status = "pending"
	statusProcessing status = "processing"
	statusRetry status = "retry"
	statusSucceeded status = "succeeded"
	statusDead status = "dead"
)

type StoredEvent struct {
	ID           string
	EventType    string
	EventPayload *bytes.Buffer
	OccurredAt   time.Time
	ProcessedAt  time.Time
}

func (e *StoredEvent) Process() {
	if e.ProcessedAt.IsZero() {
		e.ProcessedAt = time.Now()
	}
}

func RowToStoredEvent(event db.Outbox) *StoredEvent {
	return &StoredEvent{
		ID:           event.ID.String(),
		EventType:    event.EventType,
		EventPayload: bytes.NewBuffer(event.EventPayload),
		OccurredAt:   event.OccurredAt.Time,
		ProcessedAt:  event.OccurredAt.Time,
	}
}

func StoredEventToRow(event StoredEvent) db.Outbox {
	return db.Outbox{
		ID:           mapping.PgUUID(event.ID),
		EventType:    event.EventType,
		EventPayload: event.EventPayload.AvailableBuffer(),
		OccurredAt:   mapping.ToDBTimestamp(event.OccurredAt),
		ClaimedAt:  mapping.ToDBTimestamp(event.ProcessedAt),
	}
}
