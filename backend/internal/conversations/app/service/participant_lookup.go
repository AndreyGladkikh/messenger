package service

import (
	"context"
	"uuid"
)

type ParticipantLookup interface {
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}
