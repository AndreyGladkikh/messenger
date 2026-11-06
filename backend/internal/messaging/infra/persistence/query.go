package persistence

import (
	"context"
	"messenger/messenger/internal/messaging/app/query"
	"messenger/messenger/internal/platform/store"
	"uuid"
)

type QueryService struct {
	*store.Store
}

func NewQueryService(
	store *store.Store,
) *QueryService {
	return &QueryService{
		Store: store,
	}
}

func (q *QueryService) ListConversationsForUser(ctx context.Context, userID uuid.UUID) ([]query.ConversationSummary, error) {
	cnvs, err := q.Queries(ctx).ListConversationsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	cnvSummaries := make([]query.ConversationSummary, 0, len(cnvs))
	for _, cnv := range cnvs {
		cnvSummaries = append(cnvSummaries, query.ConversationSummary{
			ID:        cnv.ID,
			Type:      cnv.Type,
			Title:     cnv.Title,
			AvatarURL: cnv.AvatarUrl,
			DirectParticipantID: cnv.DirectParticipantID,
			LastMessage: query.LastMessage{
				ID:       cnv.MsgID,
				Body:     cnv.MsgBody,
				SenderID: cnv.MsgSenderID,
				SentAt:   cnv.MsgCreatedAt,
			},
		})
	}

	return cnvSummaries, nil
}
