package chat

import (
	"errors"
	"fmt"
	"messenger/messenger/internal/domain"
)

var (
	ErrDeleted           = errors.New("chat deleted")
	ErrNotFound          = fmt.Errorf("chat: %w", domain.ErrNotFound)
	ErrPrivateChatExists = errors.New("private chat already already exists")
)
