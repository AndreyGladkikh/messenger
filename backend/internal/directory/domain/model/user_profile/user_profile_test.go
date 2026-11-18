package user_profile

import (
	"testing"
	"time"

	"uuid"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	displayName := "Sponge Bob"
	avatarID := uuid.New()
	bio := "Square pants."
	birthday := time.Date(1990, time.May, 7, 0, 0, 0, 0, time.UTC)

	profile := New(
		userID,
		displayName,
		WithAvatarID(&avatarID),
		WithBio(&bio),
		WithBirthday(&birthday),
	)

	require.Equal(t, userID, profile.ID())
	require.Equal(t, displayName, profile.DisplayName)
	require.Equal(t, &avatarID, profile.AvatarID)
	require.Equal(t, &bio, profile.Bio)
	require.Equal(t, &birthday, profile.Birthday)

	events := profile.PullEvents()
	require.Len(t, events, 1)
	require.Equal(t, Created{
		UserID:      userID,
		DisplayName: displayName,
		AvatarID:    &avatarID,
		Bio:         &bio,
		Birthday:    &birthday,
	}, events[0])
	require.Equal(t, CreatedEventName, events[0].Name())
	require.Empty(t, profile.PullEvents())
}

func TestRehydrate(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	displayName := "Sponge Bob"
	profile := Rehydrate(userID, displayName, nil, nil, nil)

	require.Equal(t, userID, profile.ID())
	require.Equal(t, displayName, profile.DisplayName)
	require.Nil(t, profile.AvatarID)
	require.Nil(t, profile.Bio)
	require.Nil(t, profile.Birthday)
	require.Empty(t, profile.PullEvents())
}
