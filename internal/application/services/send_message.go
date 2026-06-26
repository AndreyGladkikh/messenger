package services

import (
	"messenger/messenger/internal/domain/chat"
	"messenger/messenger/internal/domain/message"
)

type SendMessage struct {

}

func (s *SendMessage) Execute(msg *message.Message, ch *chat.Chat) error {
	return nil
}
