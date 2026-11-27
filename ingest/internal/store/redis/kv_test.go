package redis

import (
	"errors"
	"testing"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AFK068/skilloom/ingest/internal/store/redis/mock"
)

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		rdb := mock.NewMockRDB(gomock.NewController(t))
		rdb.EXPECT().
			Get(t.Context(), "k").
			Return(goredis.NewStringResult("v", nil))

		c := &Client{RDB: rdb}
		got, err := c.Get(t.Context(), "k")
		require.NoError(t, err)
		assert.Equal(t, "v", got)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		rdb := mock.NewMockRDB(gomock.NewController(t))
		rdb.EXPECT().
			Get(t.Context(), "k").
			Return(goredis.NewStringResult("", goredis.Nil))

		c := &Client{RDB: rdb}
		got, err := c.Get(t.Context(), "k")
		require.ErrorIs(t, err, ErrNotFound)
		assert.Empty(t, got)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("redis down")
		rdb := mock.NewMockRDB(gomock.NewController(t))
		rdb.EXPECT().
			Get(t.Context(), "k").
			Return(goredis.NewStringResult("", wantErr))

		c := &Client{RDB: rdb}
		got, err := c.Get(t.Context(), "k")
		require.ErrorIs(t, err, wantErr)
		assert.Empty(t, got)
	})
}
