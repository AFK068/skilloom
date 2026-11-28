package auth_test

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/domain/principal"
	"github.com/AFK068/skilloom/ingest/internal/services/auth"
	"github.com/AFK068/skilloom/ingest/internal/services/auth/mock"
	"github.com/AFK068/skilloom/ingest/internal/store/redis"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	spaceID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	rawToken := "secret-token"
	sum := sha256.Sum256([]byte(rawToken))
	tokenKey := "token:" + hex.EncodeToString(sum[:])

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		rdb := mock.NewMockRedisStore(gomock.NewController(t))
		rdb.EXPECT().
			Get(t.Context(), tokenKey).
			Return(`{"user_id":"11111111-1111-1111-1111-111111111111","space_id":"22222222-2222-2222-2222-222222222222"}`, nil)

		svc := auth.New(auth.Deps{Redis: rdb, Log: zap.NewNop()})
		got, err := svc.Authenticate(t.Context(), rawToken)
		require.NoError(t, err)
		assert.Equal(t, principal.Principal{UserID: userID, SpaceID: spaceID}, got)
	})

	t.Run("empty token", func(t *testing.T) {
		t.Parallel()

		rdb := mock.NewMockRedisStore(gomock.NewController(t))
		svc := auth.New(auth.Deps{Redis: rdb, Log: zap.NewNop()})

		got, err := svc.Authenticate(t.Context(), "")
		require.ErrorIs(t, err, auth.ErrUnauthorized)
		assert.Zero(t, got)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		rdb := mock.NewMockRedisStore(gomock.NewController(t))
		rdb.EXPECT().
			Get(t.Context(), tokenKey).
			Return("", redis.ErrNotFound)

		svc := auth.New(auth.Deps{Redis: rdb, Log: zap.NewNop()})
		got, err := svc.Authenticate(t.Context(), rawToken)
		require.ErrorIs(t, err, auth.ErrUnauthorized)
		assert.Zero(t, got)
	})

	t.Run("redis error", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("redis down")
		rdb := mock.NewMockRedisStore(gomock.NewController(t))
		rdb.EXPECT().
			Get(t.Context(), tokenKey).
			Return("", wantErr)

		svc := auth.New(auth.Deps{Redis: rdb, Log: zap.NewNop()})
		got, err := svc.Authenticate(t.Context(), rawToken)
		require.ErrorIs(t, err, wantErr)
		require.NotErrorIs(t, err, auth.ErrUnauthorized)
		assert.Zero(t, got)
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		rdb := mock.NewMockRedisStore(gomock.NewController(t))
		rdb.EXPECT().
			Get(t.Context(), tokenKey).
			Return("{bad-json", nil)

		svc := auth.New(auth.Deps{Redis: rdb, Log: zap.NewNop()})
		got, err := svc.Authenticate(t.Context(), rawToken)
		require.Error(t, err)
		require.NotErrorIs(t, err, auth.ErrUnauthorized)
		assert.Zero(t, got)
	})

	t.Run("missing user_id", func(t *testing.T) {
		t.Parallel()

		rdb := mock.NewMockRedisStore(gomock.NewController(t))
		rdb.EXPECT().
			Get(t.Context(), tokenKey).
			Return(`{"space_id":"22222222-2222-2222-2222-222222222222"}`, nil)

		svc := auth.New(auth.Deps{Redis: rdb, Log: zap.NewNop()})
		got, err := svc.Authenticate(t.Context(), rawToken)
		require.Error(t, err)
		require.NotErrorIs(t, err, auth.ErrUnauthorized)
		assert.Zero(t, got)
	})

	t.Run("missing space_id", func(t *testing.T) {
		t.Parallel()

		rdb := mock.NewMockRedisStore(gomock.NewController(t))
		rdb.EXPECT().
			Get(t.Context(), tokenKey).
			Return(`{"user_id":"11111111-1111-1111-1111-111111111111"}`, nil)

		svc := auth.New(auth.Deps{Redis: rdb, Log: zap.NewNop()})
		got, err := svc.Authenticate(t.Context(), rawToken)
		require.Error(t, err)
		require.NotErrorIs(t, err, auth.ErrUnauthorized)
		assert.Zero(t, got)
	})
}
