package media

import (
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusReady   Status = "ready"
)

var _ domain.Aggregate = (*Media)(nil)

type Media struct {
	domain.AggregateRoot

	id     uuid.UUID
	Name   string
	Status Status
}

func New(
	name string,
) *Media {
	return &Media{
		id:   uuid.New(),
		Name: name,
		Status: StatusPending,
	}
}

func Rehydrate(
	id uuid.UUID,
	name string,
	status Status,
) *Media {
	return &Media{
		id:   id,
		Name: name,
		Status: status,
	}
}

func (m *Media) ID() uuid.UUID {
	return m.id
}

func (m *Media) MarkReady() {
	m.Status = StatusReady
}
