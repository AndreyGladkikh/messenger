package event

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/infrastructure/db"
	"messenger/messenger/internal/infrastructure/logger"
	"time"
)

type ProcessorConfig struct {
}

type Processor struct {
	eventStorage *EventStorage
	txManager    *transaction.Manager
	logger       *logger.Logger
	qs           *db.Queries
	eventBus     *Bus
	cfg          *ProcessorConfig
}

func NewProcessor(
	eventStorage *EventStorage,
	txManager *transaction.Manager,
	logger *logger.Logger,
	qs *db.Queries,
	eventBus *Bus,
) *Processor {
	return &Processor{
		eventStorage: eventStorage,
		txManager:    txManager,
		logger:       logger,
		qs:           qs,
		eventBus:     eventBus,
	}
}

func (p *Processor) Run(ctx context.Context) error {
	for {
		outboxEvent, err := p.eventStorage.GetNextUnprocessedEvent(ctx)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("retreive event to process: %w", err)
		}
		if errors.Is(err, sql.ErrNoRows) {
			time.Sleep(1 * time.Second)
			continue
		}

		domainEvent, err := translateStoredEventToDispatchedEvent(outboxEvent)
		if err != nil {
			p.eventStorage.UpdateEvent(ctx, outboxEvent, statusDead, "failed to translate stored event to dispatched event", time.Time{})
			continue
		}

		envelope := NewEnvelope(
			outboxEvent.ID.String(),
			domainEvent,
		)

		err = p.eventBus.Dispatch(ctx, envelope)
		if err != nil {
			p.eventStorage.UpdateEvent(ctx, outboxEvent, statusRetry, fmt.Sprintf("failed to process event: %s", err), time.Now().Add(5*time.Second))
			continue
		}

		p.eventStorage.UpdateEvent(ctx, outboxEvent, statusSucceeded, "", time.Time{})
	}
}