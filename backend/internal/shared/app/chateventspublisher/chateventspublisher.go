package chateventspublisher

import (
	"context"
	"messenger/messenger/internal/shared/app/event"
	"messenger/messenger/internal/shared/domain"

	"uuid"
)

type ChatEventsPublisher interface {
	Publish(ctx context.Context, chatID uuid.UUID, e event.Envelope[domain.Event]) error
}
