package group_conversation

import (
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

type GroupConversation struct {
	domain.AggregateRoot

	id   uuid.UUID
	Title string
	AvatarURL *string
}

func New(
	title string,
	avatarUrl *string,
) *GroupConversation {
	cnv := &GroupConversation{
		id:   uuid.New(),
		Title: title,
		AvatarURL: avatarUrl,
	}

	cnv.AddEvent(Created{
		ID: cnv.id,
		Title: cnv.Title,
		AvatarURL: cnv.AvatarURL,
	})

	return cnv
}

func Rehydrate(
	id uuid.UUID,
	title string,
	avatarURL *string,
) *GroupConversation {
	return &GroupConversation{
		id:   id,
		Title: title,
		AvatarURL: avatarURL,
	}
}

func (s *GroupConversation) ID() uuid.UUID {
	return s.id
}
