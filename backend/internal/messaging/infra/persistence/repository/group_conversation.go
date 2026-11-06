package repository

import (
	"context"
	"database/sql"
	"errors"
	"messenger/messenger/internal/messaging/domain/model/group_conversation"
	"messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/platform/sqlc"
	"messenger/messenger/internal/shared/app/apperr"
	"uuid"
)

var _ group_conversation.Repository = (*GroupConversationRepository)(nil)

type GroupConversationRepository struct {
	*repository.Repository
}

func NewGroupConversationRepository(repo *repository.Repository) *GroupConversationRepository {
	return &GroupConversationRepository{
		Repository: repo,
	}
}

func (c *GroupConversationRepository) Save(ctx context.Context, cnv *group_conversation.GroupConversation) error {
	err := c.Queries(ctx).UpsertGroupConversation(ctx, sqlc.UpsertGroupConversationParams{
		ConversationID: cnv.ID(),
		Title:          cnv.Title,
		AvatarUrl:      cnv.AvatarURL,
	})
	if err != nil {
		return err
	}

	c.RegisterAggregate(ctx, cnv)
	return nil
}

func (c *GroupConversationRepository) Get(ctx context.Context, id uuid.UUID) (*group_conversation.GroupConversation, error) {
	cnv, err := c.Queries(ctx).GetGroupConversation(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.NewNotFound("group conversation not found")
	}

	return group_conversation.Rehydrate(
		cnv.ConversationID,
		cnv.Title,
		cnv.AvatarUrl,
	), nil
}
