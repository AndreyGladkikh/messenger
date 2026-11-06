package messaging

import (
	"context"
	"fmt"
	"messenger/messenger/internal/bff/web/app/port"
	cnvQry "messenger/messenger/internal/messaging/app/query"
	appQuery "messenger/messenger/internal/shared/app/query"
	"uuid"
)

var _ port.MessagingReader = (*Service)(nil)

type Service struct {
	listUserConversations appQuery.Handler[cnvQry.ListUserConversations, []cnvQry.ConversationSummary]
}

func NewService(
	listUserConversations appQuery.Handler[cnvQry.ListUserConversations, []cnvQry.ConversationSummary],
) *Service {
	return &Service{
		listUserConversations: listUserConversations,
	}
}

func (s *Service) ListUserConvesations(ctx context.Context, userID uuid.UUID) ([]port.ConversationSummary, error) {
	cnvs, err := s.listUserConversations.Handle(ctx, cnvQry.ListUserConversations{
		UserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("list user conversations: %w", err)
	}

	cnvsSum := make([]port.ConversationSummary, 0, len(cnvs))
	for _, c := range cnvs {
		cnvsSum = append(cnvsSum, port.ConversationSummary{
			ID:        c.ID,
			Type:      c.Type,
			Title:     c.Title,
			AvatarURL: c.AvatarURL,
			LastMessage: port.LastMessage{
				ID:       c.LastMessage.ID,
				Body:     c.LastMessage.Body,
				SenderID: c.LastMessage.SenderID,
				SentAt:   c.LastMessage.SentAt,
			},
		})
	}

	return cnvsSum, nil
}
