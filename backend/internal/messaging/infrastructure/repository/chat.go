package repository

import (
	"context"
	"database/sql"
	"errors"
	"messenger/messenger/internal/messaging/domain/chat"
	"messenger/messenger/internal/messaging/infrastructure/sqlc"
	"messenger/messenger/internal/platform/db"
	"messenger/messenger/internal/platform/postgres"
	"messenger/messenger/internal/platform/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type ChatRepository struct {
	*repository.Repository
}

func NewChatRepository(
	repo *repository.Repository,
) *ChatRepository {
	return &ChatRepository{
		Repository: repo,
	}
}

func (r *ChatRepository) Add(ctx context.Context, chat *chat.Chat) error {
	err := r.Queries(ctx).CreateChat(ctx, sqlc.CreateChatParams{
		ID:   chat.ID(),
		Kind: string(chat.Type()),
		Name: db.ToDBText(chat.Name()),
	})
	if err != nil {
		return err
	}

	r.RegisterAggregate(ctx, chat)

	return err
}

func (r *ChatRepository) AddPrivate(ctx context.Context, c *chat.Chat, privateChatPair *chat.PrivateChatPair) error {
	err := r.Add(ctx, c)
	if err != nil {
		return err
	}

	err = r.Queries(ctx).CreatePrivateChat(ctx, sqlc.CreatePrivateChatParams{
		ChatID:              c.ID(),
		FirstParticipantID:  privateChatPair.FirstParticipantID(),
		SecondParticipantID: privateChatPair.SecondParticipantID(),
	})
	if err != nil {
		if err, ok := errors.AsType[*pgconn.PgError](err); ok && err.Code == postgres.UniqueViolationErrCode && err.ConstraintName == "private_chats_first_participant_id_second_participant_id_key" {
			return chat.ErrPrivateChatAlreadyExists
		}
		return err
	}

	return err
}

func (r *ChatRepository) PrivateChatExists(ctx context.Context, participant1, participant2 uuid.UUID) (bool, error) {
	exists, err := r.Queries(ctx).PrivateChatExists(ctx, sqlc.PrivateChatExistsParams{
		ParticipantID:   participant1,
		ParticipantID_2: participant2,
	})
	if err != nil {
		return false, err
	}

	return bool(exists), nil
}

func (r *ChatRepository) Get(ctx context.Context, chatID uuid.UUID) (*chat.Chat, error) {
	chatRow, err := r.Queries(ctx).GetChat(ctx, chatID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, chat.ErrNotFound
	}

	chat := chat.Rehydrate(
		chatRow.ID,
		chat.ChatType(chatRow.Kind),
		chatRow.Name.String,
	)

	return chat, nil
}
