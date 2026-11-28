package messages_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/domain/message"
	"github.com/AFK068/skilloom/ingest/internal/domain/principal"
	"github.com/AFK068/skilloom/ingest/internal/services/messages"
	"github.com/AFK068/skilloom/ingest/internal/services/messages/mock"
)

func TestSaveMessages(t *testing.T) {
	t.Parallel()

	p := principal.Principal{
		UserID:  uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		SpaceID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
	}
	capturedAt := time.Date(2026, 11, 28, 12, 0, 0, 0, time.UTC)

	msg1 := message.Message{
		EventID:    uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		Harness:    "claude",
		SessionID:  "sess-1",
		Role:       "user",
		Text:       "hello",
		CapturedAt: capturedAt,
	}
	msg2 := message.Message{
		EventID:    uuid.MustParse("44444444-4444-4444-4444-444444444444"),
		Harness:    "claude",
		SessionID:  "sess-1",
		Role:       "assistant",
		Text:       "hi",
		CapturedAt: capturedAt.Add(time.Second),
	}

	t.Run("ok batch", func(t *testing.T) {
		t.Parallel()

		pub := mock.NewMockPublisher(gomock.NewController(t))
		gomock.InOrder(
			pub.EXPECT().PublishMessage(t.Context(), p, msg1).Return(nil),
			pub.EXPECT().PublishMessage(t.Context(), p, msg2).Return(nil),
		)

		svc := messages.New(messages.Deps{Publisher: pub, Log: zap.NewNop()})
		require.NoError(t, svc.SaveMessages(t.Context(), p, []message.Message{msg1, msg2}))
	})

	t.Run("empty batch", func(t *testing.T) {
		t.Parallel()

		pub := mock.NewMockPublisher(gomock.NewController(t))
		svc := messages.New(messages.Deps{Publisher: pub, Log: zap.NewNop()})

		require.ErrorIs(t, svc.SaveMessages(t.Context(), p, nil), messages.ErrInvalidMessage)
	})

	t.Run("invalid message skips publishing whole batch", func(t *testing.T) {
		t.Parallel()

		invalid := msg2
		invalid.EventID = uuid.Nil

		pub := mock.NewMockPublisher(gomock.NewController(t))
		svc := messages.New(messages.Deps{Publisher: pub, Log: zap.NewNop()})

		require.ErrorIs(t, svc.SaveMessages(t.Context(), p, []message.Message{msg1, invalid}), messages.ErrInvalidMessage)
	})

	t.Run("publish error", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("nats down")
		pub := mock.NewMockPublisher(gomock.NewController(t))
		pub.EXPECT().PublishMessage(t.Context(), p, msg1).Return(wantErr)

		svc := messages.New(messages.Deps{Publisher: pub, Log: zap.NewNop()})
		err := svc.SaveMessages(t.Context(), p, []message.Message{msg1})
		require.ErrorIs(t, err, wantErr)
		require.NotErrorIs(t, err, messages.ErrInvalidMessage)
	})
}
