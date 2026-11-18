package handle

import (
	"messenger/messenger/internal/directory/domain/model"
	"messenger/messenger/internal/shared/domain"
)

const CreatedEventName = "directory.handle.Created.v1"

var _ domain.Event = Created{}

type Created struct {
	Target model.TargetRef
	Value      string
}

func (c Created) Name() string {
	return CreatedEventName
}
