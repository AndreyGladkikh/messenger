package query

import "context"

type Reader interface {
	Search(ctx context.Context, name string)
}
