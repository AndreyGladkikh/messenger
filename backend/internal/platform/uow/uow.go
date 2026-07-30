package uow

import (
	"maps"
	"messenger/messenger/internal/messaging/domain"
	"slices"
)

type UnitOfWork struct {
	aggregates map[string]domain.Aggregate
}

func New() *UnitOfWork {
	return &UnitOfWork{
		aggregates: make(map[string]domain.Aggregate),
	}
}

func (u *UnitOfWork) RegisterAggregate(a domain.Aggregate) {
	u.aggregates[a.ID().String()] = a
}

func (u *UnitOfWork) Aggregates() []domain.Aggregate {
	return slices.Collect(maps.Values(u.aggregates))
}
