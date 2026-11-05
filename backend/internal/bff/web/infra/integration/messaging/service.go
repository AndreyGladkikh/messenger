package messaging

import (
	"context"
	"fmt"
	"messenger/messenger/internal/bff/web/application/port"
	"messenger/messenger/internal/conversations/app/query"
	"messenger/messenger/internal/platform/querybus"
	"uuid"
)

var _ port.MessagingReader = (*Service)(nil)

type Service struct {
	queryBus *querybus.Bus
}

func NewService(
	queryBus *querybus.Bus,
) *Service {
	return &Service{
		queryBus: queryBus,
	}
}

func (s *Service) ListUserConvesations(ctx context.Context, userID uuid.UUID) ([]port.ConversationSummary, error) {
	cnvs, err := s.queryBus.Dispatch[[]query.ConversationSummary](ctx, query.ListUserConversations{
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
