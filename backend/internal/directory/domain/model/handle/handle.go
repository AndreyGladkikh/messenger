package handle

import (
	"messenger/messenger/internal/directory/domain/model"
	"messenger/messenger/internal/shared/domain"

	"uuid"
)

var _ domain.Aggregate = (*Handle)(nil)

type Handle struct {
	domain.AggregateRoot

	Target model.TargetRef
	Value  string
}

func New(
	target model.TargetRef,
	value string,
) *Handle {
	h := &Handle{
		Target: target,
		Value:  value,
	}
	h.AddEvent(Created{
		Target: target,
		Value:  h.Value,
	})

	return h
}

func Rehydrate(
	target model.TargetRef,
	value string,
) *Handle {
	return &Handle{
		Target: target,
		Value:  value,
	}
}

func (h *Handle) ID() uuid.UUID {
	return h.Target.ID
}
