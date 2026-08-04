package uow

import (
	"maps"
	sharedDomain "messenger/messenger/internal/shared/domain"
	"slices"
)

type UnitOfWork struct {
	aggregates map[string]sharedDomain.Aggregate
}

func New() *UnitOfWork {
	return &UnitOfWork{
		aggregates: make(map[string]sharedDomain.Aggregate),
	}
}

func (u *UnitOfWork) RegisterAggregate(a sharedDomain.Aggregate) {
	u.aggregates[a.ID().String()] = a
}

func (u *UnitOfWork) Aggregates() []sharedDomain.Aggregate {
	return slices.Collect(maps.Values(u.aggregates))
}
