package ws_test

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/core/mock"
	"github.com/TerraDharitri/drt-go-chain-notifier/common"
	"github.com/TerraDharitri/drt-go-chain-notifier/dispatcher/ws"
	"github.com/TerraDharitri/drt-go-chain-notifier/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockArgsWSHandler() ws.ArgsWebSocketProcessor {
	return ws.ArgsWebSocketProcessor{
		Dispatcher: &mocks.HubStub{},
		Upgrader:   &mocks.WSUpgraderStub{},
		Marshaller: &mock.MarshalizerMock{},
	}
}

func TestNewWebSocketHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil dispatcher handler", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWSHandler()
		args.Dispatcher = nil

		wh, err := ws.NewWebSocketProcessor(args)
		require.True(t, check.IfNil(wh))
		assert.Equal(t, ws.ErrNilDispatcher, err)
	})

	t.Run("nil ws upgrader", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWSHandler()
		args.Upgrader = nil

		wh, err := ws.NewWebSocketProcessor(args)
		require.True(t, check.IfNil(wh))
		assert.Equal(t, ws.ErrNilWSUpgrader, err)
	})

	t.Run("nil marshaller", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWSHandler()
		args.Marshaller = nil

		wh, err := ws.NewWebSocketProcessor(args)
		require.True(t, check.IfNil(wh))
		assert.Equal(t, common.ErrNilMarshaller, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWSHandler()
		wh, err := ws.NewWebSocketProcessor(args)
		require.False(t, check.IfNil(wh))
		require.Nil(t, err)
	})
}
