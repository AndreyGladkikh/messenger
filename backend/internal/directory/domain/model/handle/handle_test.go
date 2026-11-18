package handle

import (
	"testing"

	"messenger/messenger/internal/directory/domain/model"

	"uuid"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	target := model.TargetRef{
		ID: uuid.New(),
		Type: model.TargetTypeUser,
	}
	value := "MixedCaseHandle"

	handle := New(target, value)

	require.Equal(t, target.ID, handle.ID())
	require.Equal(t, target.ID, handle.Target.ID)
	require.Equal(t, target.Type, handle.Target.Type)
	require.Equal(t, value, handle.Value)

	events := handle.PullEvents()
	require.Len(t, events, 1)
	require.Equal(t, Created{
		Target: target,
		Value:      value,
	}, events[0])
	require.Equal(t, CreatedEventName, events[0].Name())
	require.Empty(t, handle.PullEvents())
}

func TestRehydrate(t *testing.T) {
	t.Parallel()

	target := model.TargetRef{
		ID: uuid.New(),
		Type: model.TargetTypeBot,
	}
	value := "BotHandle"

	handle := Rehydrate(target, value)

	require.Equal(t, target.ID, handle.ID())
	require.Equal(t, target.ID, handle.Target.ID)
	require.Equal(t, target.Type, handle.Target.Type)
	require.Equal(t, value, handle.Value)
	require.Empty(t, handle.PullEvents())
}
