package user_profile

import (
	"messenger/messenger/internal/shared/domain"
	"time"

	"uuid"
)

var _ domain.Aggregate = (*UserProfile)(nil)

type UserProfile struct {
	domain.AggregateRoot

	UserID      uuid.UUID
	DisplayName string
	AvatarID    *uuid.UUID
	Bio         *string
	Birthday    *time.Time
}

func New(
	userID uuid.UUID,
	displayName string,
	opts ...Option,
) *UserProfile {
	p := &UserProfile{
		UserID:      userID,
		DisplayName: displayName,
	}
	for _, opt := range opts {
		opt(p)
	}

	p.AddEvent(Created{
		UserID:      p.UserID,
		DisplayName: p.DisplayName,
		AvatarID:    p.AvatarID,
		Bio:         p.Bio,
		Birthday:    p.Birthday,
	})

	return p
}

func (p *UserProfile) Update(opts ...Option) {
	for _, opt := range opts {
		opt(p)
	}

	p.AddEvent(Updated{
		UserID:      p.UserID,
		DisplayName: p.DisplayName,
		AvatarID:    p.AvatarID,
		Bio:         p.Bio,
		Birthday:    p.Birthday,
	})
}

func Rehydrate(
	userID uuid.UUID,
	displayName string,
	avatarID *uuid.UUID,
	bio *string,
	birthday *time.Time,
) *UserProfile {
	return &UserProfile{
		UserID:      userID,
		DisplayName: displayName,
		AvatarID:    avatarID,
		Bio:         bio,
		Birthday:    birthday,
	}
}

func (p *UserProfile) ID() uuid.UUID {
	return p.UserID
}

type Option func(*UserProfile)

func WithDisplayName(dispName string) Option {
	return func(p *UserProfile) {
		p.DisplayName = dispName
	}
}

func WithAvatarID(avatarID *uuid.UUID) Option {
	return func(p *UserProfile) {
		p.AvatarID = avatarID
	}
}

func WithBio(bio *string) Option {
	return func(p *UserProfile) {
		p.Bio = bio
	}
}

func WithBirthday(birthday *time.Time) Option {
	return func(p *UserProfile) {
		p.Birthday = birthday
	}
}
