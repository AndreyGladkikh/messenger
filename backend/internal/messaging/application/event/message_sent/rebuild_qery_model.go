package message_sent

import (
	"context"
	"fmt"
	"messenger/messenger/internal/messaging/domain/message"
	"messenger/messenger/internal/shared/application/event"
)

type RebuildQueryModelHandler struct {
}

func NewRebuildQueryModelHandler() *RebuildQueryModelHandler {
	return &RebuildQueryModelHandler{}
}

func (h *RebuildQueryModelHandler) Handle(ctx context.Context, e event.Envelope[message.MessageSent]) error {
	fmt.Printf("event handler %s executed", h.Name())
	return fmt.Errorf("false error")
	// return nil
}

func (h *RebuildQueryModelHandler) Name() string {
	return "event_handler.rebuild_query_model.v1"
}
