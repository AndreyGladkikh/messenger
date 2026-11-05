package realtime

import (
	"context"
	"encoding/json/v2"
	"errors"
	cnvEvent "messenger/messenger/internal/conversations/app/event"
	drctrQry "messenger/messenger/internal/directory/app/query"
	"messenger/messenger/internal/directory/domain/entry"
	"messenger/messenger/internal/platform/rabbitmq"
	appEvent "messenger/messenger/internal/shared/app/event"
	appQuery "messenger/messenger/internal/shared/app/query"
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersonalDeliveryHandlerPreservesGetEntryError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("directory unavailable")
	senderID := uuid.New()
	handler := appQuery.HandlerFunc[drctrQry.GetEntry, drctrQry.EntryView](
		func(_ context.Context, qry drctrQry.GetEntry) (drctrQry.EntryView, error) {
			assert.Equal(t, senderID, qry.ID)
			assert.Equal(t, entry.TargetTypeUser.String(), qry.Type)
			return drctrQry.EntryView{}, expectedErr
		},
	)
	payload, err := json.Marshal(appEvent.Envelope[cnvEvent.DirectConversationStarted]{
		ID:         uuid.New(),
		OccurredAt: time.Now(),
		Event: cnvEvent.DirectConversationStarted{
			ConversationID: uuid.New(),
			SenderID:       senderID,
			RecipientID:    uuid.New(),
			MessageID:      uuid.New(),
			MessageBody:    "Hello",
		},
	})
	require.NoError(t, err)
	consumer := NewPersonalDeliveryHandler(nil, handler)

	err = consumer.Handle(t.Context(), rabbitmq.Envelope{
		Type:    cnvEvent.DirectConversationStartedEventName,
		Payload: payload,
	})

	require.ErrorIs(t, err, expectedErr)
}
