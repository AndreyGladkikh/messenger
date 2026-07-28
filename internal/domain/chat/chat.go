package chat

import (
	"messenger/messenger/internal/domain"
	"time"
)

type Chat struct {
	domain.BaseAggregate
	id        string
	typ       ChatType
	name      string
	deletedAt time.Time
}

func Create(
	id string,
	typ ChatType,
) *Chat {
	c := new(Chat)
	c.id = id
	c.typ = typ

	c.AddEvent(ChatCreated{
		ChatID: id,
	})

	return c
}

func Rehydrate(
	id string,
	typ ChatType,
	name string,
) *Chat {
	return &Chat{
		id:   id,
		typ:  typ,
		name: name,
	}
}

func (c *Chat) ID() string {
	return c.id
}

func (c *Chat) SetID(id string) {
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
