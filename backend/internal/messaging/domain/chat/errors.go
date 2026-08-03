package chat

import (
	"errors"
	"fmt"
	sharedDomain "messenger/messenger/internal/messaging/domain"
)

var (
	ErrDeleted                  = errors.New("chat deleted")
	ErrNotFound                 = fmt.Errorf("chat: %w", sharedDomain.ErrNotFound)
	ErrPrivateChatAlreadyExists = errors.New("private chat already already exists")
)
