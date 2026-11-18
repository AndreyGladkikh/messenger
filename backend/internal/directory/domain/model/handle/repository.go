package handle

import (
	"context"
	"messenger/messenger/internal/directory/domain/model"
)

type Repository interface {
	Save(context.Context, *Handle) error
	DeleteByTarget(context.Context, model.TargetRef) error
}
