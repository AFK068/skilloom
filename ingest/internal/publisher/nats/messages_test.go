package natspub_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/domain/message"
	"github.com/AFK068/skilloom/ingest/internal/domain/principal"
	natspub "github.com/AFK068/skilloom/ingest/internal/publisher/nats"
	"github.com/AFK068/skilloom/ingest/internal/publisher/nats/mock"
)

func TestPublishMessage(t *testing.T) {
	t.Parallel()

	p := principal.Principal{
		UserID:  uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		SpaceID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
	}
	eventID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	capturedAt := time.Date(2026, 11, 28, 12, 0, 0, 0, time.UTC)

	msg := message.Message{
		EventID:    eventID,
		Harness:    "claude",
		SessionID:  "sess-1",
		TurnID:     "turn-1",
		Role:       "user",
		Text:       "hello world",
		CapturedAt: capturedAt,
		Meta:       json.RawMessage(`{"source":"cli"}`),
	}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		client := mock.NewMockClient(gomock.NewController(t))
		var published []byte
		client.EXPECT().
			Publish(t.Context(), "ingest.messages", eventID.String(), gomock.Any()).
			DoAndReturn(func(_ any, _, _ string, data []byte) error {
				published = data
				return nil
			})

		pub := natspub.New(natspub.Deps{Client: client, Log: zap.NewNop()})
		require.NoError(t, pub.PublishMessage(t.Context(), p, msg))

		assert.JSONEq(t, `{
			"event_id": "33333333-3333-3333-3333-333333333333",
			"user_id": "11111111-1111-1111-1111-111111111111",
			"space_id": "22222222-2222-2222-2222-222222222222",
			"harness": "claude",
			"session_id": "sess-1",
			"turn_id": "turn-1",
			"role": "user",
			"text": "hello world",
			"captured_at": "2026-11-28T12:00:00Z",
			"meta": {"source": "cli"}
		}`, string(published))
	})

	t.Run("ok omits empty optional fields", func(t *testing.T) {
		t.Parallel()

		minimal := msg
		minimal.TurnID = ""
		minimal.Meta = nil

		client := mock.NewMockClient(gomock.NewController(t))
		var published []byte
		client.EXPECT().
			Publish(t.Context(), "ingest.messages", eventID.String(), gomock.Any()).
			DoAndReturn(func(_ any, _, _ string, data []byte) error {
				published = data
				return nil
			})

		pub := natspub.New(natspub.Deps{Client: client, Log: zap.NewNop()})
		require.NoError(t, pub.PublishMessage(t.Context(), p, minimal))

		assert.JSONEq(t, `{
			"event_id": "33333333-3333-3333-3333-333333333333",
			"user_id": "11111111-1111-1111-1111-111111111111",
			"space_id": "22222222-2222-2222-2222-222222222222",
			"harness": "claude",
			"session_id": "sess-1",
			"role": "user",
			"text": "hello world",
			"captured_at": "2026-11-28T12:00:00Z"
		}`, string(published))
	})

	t.Run("publish error", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("nats down")
		client := mock.NewMockClient(gomock.NewController(t))
		client.EXPECT().
			Publish(t.Context(), "ingest.messages", eventID.String(), gomock.Any()).
			Return(wantErr)

		pub := natspub.New(natspub.Deps{Client: client, Log: zap.NewNop()})
		require.ErrorIs(t, pub.PublishMessage(t.Context(), p, msg), wantErr)
	})
}
