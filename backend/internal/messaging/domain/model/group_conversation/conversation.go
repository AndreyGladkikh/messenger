package group_conversation

import (
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

type GroupConversation struct {
	domain.AggregateRoot

	id       uuid.UUID
	Title    string
	AvatarID *uuid.UUID
}

func New(
	title string,
	avatarID *uuid.UUID,
) *GroupConversation {
	cnv := &GroupConversation{
		id:       uuid.New(),
		Title:    title,
		AvatarID: avatarID,
	}

	cnv.AddEvent(Created{
		ID:       cnv.id,
		Title:    cnv.Title,
		AvatarID: cnv.AvatarID,
	})

	return cnv
}

func Rehydrate(
	id uuid.UUID,
	title string,
	avatarID *uuid.UUID,
) *GroupConversation {
	return &GroupConversation{
		id:       id,
		Title:    title,
		AvatarID: avatarID,
	}
}

func (s *GroupConversation) ID() uuid.UUID {
	return s.id
}
