package query

import "context"

type GlobalSearch struct {
	Str string
}

type GlobalSearchHandler struct {

}

func (h *GlobalSearchHandler) Handle(ctx context.Context, qry GlobalSearch) {
	
}
