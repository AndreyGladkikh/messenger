package chat

import (
	"messenger/messenger/internal/messaging/domain"
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	domain.BaseAggregate
	id        uuid.UUID
	typ       ChatType
	name      string
	deletedAt time.Time
}

func Create(
	typ ChatType,
) *Chat {
	chat := &Chat{
		id:  uuid.New(),
		typ: typ,
	}

	chat.AddEvent(ChatCreated{
		ChatID: chat.ID(),
	})

	return chat
}

func Rehydrate(
	id uuid.UUID,
	typ ChatType,
	name string,
) *Chat {
	return &Chat{
		id:   id,
		typ:  typ,
		name: name,
	}
}

func (c *Chat) ID() uuid.UUID {
	return c.id
}

func (c *Chat) SetID(id uuid.UUID) {
	c.id = id
}

func (c *Chat) Type() ChatType {
	return c.typ
}

func (c *Chat) SetType(typ ChatType) {
	c.typ = typ
}

func (c *Chat) Name() string {
	return c.name
}

func (c *Chat) SetName(name string) {
	c.name = name
}

func (c *Chat) IsDeleted() bool {
	return !c.deletedAt.IsZero()
}
