package uow

import (
	"context"
	"maps"
	"messenger/messenger/internal/shared/app/event"
	"messenger/messenger/internal/shared/domain"
	"slices"
	"sync"
)

type UnitOfWork struct {
	aggregatesMu sync.Mutex
	aggregates map[string]domain.Aggregate
	hooksMu sync.Mutex
	hooks []Hook
	eventsMu sync.Mutex
	events []event.Event
}

func New() *UnitOfWork {
	return &UnitOfWork{
		aggregates: make(map[string]domain.Aggregate),
	}
}

func (u *UnitOfWork) RegisterAggregate(a domain.Aggregate) {
	u.aggregatesMu.Lock()
	defer u.aggregatesMu.Unlock()
	u.aggregates[a.ID().String()] = a
}

func (u *UnitOfWork) Aggregates() []domain.Aggregate {
	u.aggregatesMu.Lock()
	defer u.aggregatesMu.Unlock()
	return slices.Collect(maps.Values(u.aggregates))
}

func (u *UnitOfWork) AddHook(h Hook) {
	u.hooksMu.Lock()
	defer u.hooksMu.Unlock()
	u.hooks = append(u.hooks, h)
}

func (u *UnitOfWork) Hooks() []Hook {
	u.hooksMu.Lock()
	defer u.hooksMu.Unlock()
	return slices.Clone(u.hooks)
}

func (u *UnitOfWork) AddEvent(e event.Event) {
	u.eventsMu.Lock()
	defer u.eventsMu.Unlock()
	u.events = append(u.events, e)
}

func (u *UnitOfWork) Events() []event.Event {
	u.eventsMu.Lock()
	defer u.eventsMu.Unlock()
	return slices.Clone(u.events)
}

type Hook func(context.Context) error
