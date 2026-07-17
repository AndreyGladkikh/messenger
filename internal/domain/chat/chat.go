package chat

import "messenger/messenger/internal/domain"

type Chat struct {
	domain.BaseAggregate
	id   string
	typ  ChatType
	name string
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
