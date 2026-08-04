package chat

import (
	"fmt"
	sharedDomain "messenger/messenger/internal/shared/domain"
)

var (
	ErrDeleted                  = fmt.Errorf("chat: %w", sharedDomain.ErrDeleted)
	ErrNotFound                 = fmt.Errorf("chat: %w", sharedDomain.ErrNotFound)
	ErrPrivateChatAlreadyExists = fmt.Errorf("private chat: %w", sharedDomain.ErrAlreadyExists)
)
