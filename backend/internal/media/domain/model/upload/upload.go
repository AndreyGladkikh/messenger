package upload

import (
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
)

var _ domain.Aggregate = (*Upload)(nil)

type Upload struct {
	domain.AggregateRoot

	id          uuid.UUID
	InitiatorID uuid.UUID
	Status      Status
	MediaID     uuid.UUID
}

func New(
	initiatorID uuid.UUID,
	mediaID uuid.UUID,
) *Upload {
	return &Upload{
		id:          uuid.New(),
		InitiatorID: initiatorID,
		Status:      StatusPending,
		MediaID:     mediaID,
	}
}

func Rehydrate(
	id uuid.UUID,
	initiatorID uuid.UUID,
	status Status,
	mediaID uuid.UUID,
) *Upload {
	return &Upload{
		id:          id,
		InitiatorID: initiatorID,
		Status:      status,
		MediaID:     mediaID,
	}
}

func (u *Upload) ID() uuid.UUID {
	return u.id
}

func (u *Upload) Complete() {
	u.Status = StatusCompleted
}
