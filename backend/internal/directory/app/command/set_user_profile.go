package command

import (
	"context"
	"errors"
	"fmt"
	"messenger/messenger/internal/directory/app/port"
	"messenger/messenger/internal/directory/domain/model"
	"messenger/messenger/internal/directory/domain/model/handle"
	"messenger/messenger/internal/directory/domain/model/user_profile"
	"messenger/messenger/internal/shared/app/apperr"
	"messenger/messenger/internal/shared/app/command"
	"time"
	"uuid"
)

const SetUserProfileCommandName = "directory.SetUserProfile.v1"

var _ command.Command = SetUserProfile{}

type SetUserProfile struct {
	UserID      uuid.UUID
	Handle      *string
	DisplayName string
	AvatarID    *uuid.UUID
	Bio         *string
	Birthday    *time.Time
}

func (r SetUserProfile) IsCommand() {}

func (r SetUserProfile) Name() string {
	return SetUserProfileCommandName
}

type SetUserProfileHandler struct {
	identityReader port.IdentityReader
	mediaReader           port.MediaReader
	userProfileRepository user_profile.Repository
	handleRepository      handle.Repository
}

func NewSetUserProfileHandler(
	identityReader port.IdentityReader,
	mediaReader           port.MediaReader,
	userProfileRepository user_profile.Repository,
	handleRepository      handle.Repository,
) *SetUserProfileHandler {
	return &SetUserProfileHandler{
		mediaReader: mediaReader,
		identityReader: identityReader,
		userProfileRepository: userProfileRepository,
		handleRepository: handleRepository,
	}
}

func (h *SetUserProfileHandler) Handle(ctx context.Context, cmd SetUserProfile) (any, error) {
	exists, err := h.identityReader.ExistsUser(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, apperr.NewNotFound(fmt.Sprintf("User with id %q doesn't exist", cmd.UserID))
	}

	isNew := false

	p, err := h.userProfileRepository.Get(ctx, cmd.UserID)
	if err != nil && !errors.Is(err, apperr.ErrNotFound) {
		return nil, err
	}

	if errors.Is(err, apperr.ErrNotFound) {
		isNew = true
	}

	if cmd.AvatarID != nil && (isNew || cmd.AvatarID != p.AvatarID) {
		exists, err := h.mediaReader.MediaExists(ctx, *cmd.AvatarID)
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, apperr.NewNotFound(fmt.Sprintf("Avatar with id %q not found", cmd.AvatarID))
		}
	}

	target := model.TargetRef{
		ID:   cmd.UserID,
		Type: model.TargetTypeUser,
	}
	if cmd.Handle != nil {
		handle := handle.New(
			target,
			*cmd.Handle,
		)

		err := h.handleRepository.Save(ctx, handle)
		if err != nil {
			return nil, err
		}
	} else {
		err := h.handleRepository.DeleteByTarget(ctx, target)
		if err != nil {
			return nil, err
		}
	}

	if isNew {
		p = user_profile.New(
			cmd.UserID,
			cmd.DisplayName,
			user_profile.WithBio(cmd.Bio),
			user_profile.WithBirthday(cmd.Birthday),
			user_profile.WithAvatarID(cmd.AvatarID),
		)
	} else {
		p.Update(
			user_profile.WithDisplayName(cmd.DisplayName),
			user_profile.WithBio(cmd.Bio),
			user_profile.WithBirthday(cmd.Birthday),
			user_profile.WithAvatarID(cmd.AvatarID),
		)
	}

	err = h.userProfileRepository.Save(ctx, p)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
