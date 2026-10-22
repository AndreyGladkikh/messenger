package command

import (
	"context"
	"messenger/messenger/internal/directory/domain/entry"
	"messenger/messenger/internal/shared/app/apperr"
	"messenger/messenger/internal/shared/app/command"
	"uuid"
)

var _ command.Command = RegisterEntry{}

const RegisterEntryCommandName = "RegisterEntry"

type RegisterEntry struct {
	TargetID uuid.UUID
	TargetType string
	DisplayName string
	Handle string
}

func (r RegisterEntry) IsCommand() {}

func (r RegisterEntry) Name() string {
	return RegisterEntryCommandName
}

var _ command.Handler[RegisterEntry, any] = (*RegisterEntryHandler)(nil)

type RegisterEntryHandler struct {
	entryRepo entry.Repository
	targetLookup *TargetLookup
}

func NewRegisterEntryHandler(
	entryRepo entry.Repository,
	targetLookup *TargetLookup,
) *RegisterEntryHandler {
	return &RegisterEntryHandler{
		entryRepo: entryRepo,
		targetLookup: targetLookup,
	}
}

func (h *RegisterEntryHandler) Handle(ctx context.Context, qry RegisterEntry) (any, error) {
	trgt, err := entry.ParseTargetType(qry.TargetType)
	if err != nil {
		return nil, err
	}

	exists, err := h.targetLookup.Exists(ctx, qry.TargetID, trgt)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, apperr.NewNotFound("Entry target doesn't exist")
	}

	entr := entry.New(
		qry.TargetID,
		trgt,
		qry.DisplayName,
		qry.Handle,
	)

	err = h.entryRepo.Save(ctx, entr)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

type UserLookup interface {
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type TargetLookup struct {
	userLookup UserLookup
}

func NewTargetLookup(
	userLookup UserLookup,
) *TargetLookup {
	return &TargetLookup{
		userLookup: userLookup,
	}
}

func (t *TargetLookup) Exists(ctx context.Context, id uuid.UUID, typ entry.TargetType) (bool, error) {
	switch typ{
	case entry.TargetTypeUser:
		return t.userLookup.Exists(ctx, id)
	default:
		panic("not implemented")
	}
}
