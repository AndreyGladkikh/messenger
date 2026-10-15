package entry

import "context"

type Repository interface {
	Save(context.Context, *Entry) error
}
