package repository

import (
	"context"
	"database/sql"
	"errors"
	"messenger/messenger/internal/conversations/domain/model/conversation"
	"messenger/messenger/internal/platform/postgres"
	"messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/platform/sqlc"
	"messenger/messenger/internal/shared/domain"
	"uuid"

	"github.com/jackc/pgx/v5/pgconn"
)

type ConversationRepository struct {
	*repository.Repository
}

func NewConversationRepository(repo *repository.Repository) *ConversationRepository {
	return &ConversationRepository{
		Repository: repo,
	}
}

func (c *ConversationRepository) GetDirectBetween(ctx context.Context, participants conversation.DirectConversationParticipantsPair) (*conversation.Conversation, error) {
	cnv, err := c.Queries(ctx).GetDirectConversationBetween(ctx, sqlc.GetDirectConversationBetweenParams{
		ParticipantIDLow:  participants.ParticipantLow(),
		ParticipantIDHigh: participants.ParticipantHigh(),
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}

	return conversation.Rehydrate(
		cnv.ID,
		conversation.ConversationType(cnv.Type),
	), nil
}

func (c *ConversationRepository) Save(ctx context.Context, cnv *conversation.Conversation) error {
	err := c.Queries(ctx).CreateConversation(ctx, sqlc.CreateConversationParams{
		ID:   cnv.ID(),
		Type: cnv.Type.String(),
	})
	if err != nil {
		return err
	}

	c.RegisterAggregate(ctx, cnv)
	return nil
}

func (c *ConversationRepository) Get(ctx context.Context, id uuid.UUID) (*conversation.Conversation, error) {
	cnv, err := c.Queries(ctx).GetConversation(ctx, id)
	if err != nil {
		return nil, err
	}

	return conversation.Rehydrate(
		cnv.ID,
		conversation.ConversationType(cnv.Type),
	), nil
}

func (c *ConversationRepository) SaveDirect(ctx context.Context, cnv *conversation.Conversation, participants conversation.DirectConversationParticipantsPair) error {
	err := c.Save(ctx, cnv)
	if err != nil {
		return err
	}

	err = c.Queries(ctx).CreateDirectConversation(ctx, sqlc.CreateDirectConversationParams{
		ConversationID:    cnv.ID(),
		ParticipantIDLow:  participants.ParticipantLow(),
		ParticipantIDHigh: participants.ParticipantHigh(),
	})
	if err != nil {
		if err, ok := errors.AsType[*pgconn.PgError](err); ok && err.Code == postgres.UniqueViolationErrCode && err.ConstraintName == "direct_conversations_conversation_id_key" {
			return domain.ErrAlreadyExists
		}
		return err
	}

	return nil
}
