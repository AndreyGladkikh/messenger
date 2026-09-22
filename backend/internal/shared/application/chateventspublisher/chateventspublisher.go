package chateventspublisher

import (
	"context"
	"messenger/messenger/internal/shared/application/event"
	"messenger/messenger/internal/shared/domain"

	"github.com/google/uuid"
)

type ChatEventsPublisher interface {
	Publish(ctx context.Context, chatID uuid.UUID, e event.Envelope[domain.Event]) error
}
