package message_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/AFK068/skilloom/ingest/internal/domain/message"
)

func TestMessageValidate(t *testing.T) {
	t.Parallel()

	valid := message.Message{
		EventID:    uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		Harness:    "claude",
		SessionID:  "sess-1",
		TurnID:     "turn-1",
		Role:       "user",
		Text:       "hello world",
		CapturedAt: time.Date(2026, 11, 28, 12, 0, 0, 0, time.UTC),
		Meta:       json.RawMessage(`{"source":"cli"}`),
	}

	t.Run("valid full", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, valid.Validate())
	})

	t.Run("valid without optional fields", func(t *testing.T) {
		t.Parallel()

		msg := valid
		msg.TurnID = ""
		msg.Meta = nil

		require.NoError(t, msg.Validate())
	})

	cases := []struct {
		name   string
		mutate func(*message.Message)
	}{
		{name: "nil event_id", mutate: func(m *message.Message) { m.EventID = uuid.Nil }},
		{name: "empty harness", mutate: func(m *message.Message) { m.Harness = "" }},
		{name: "empty session_id", mutate: func(m *message.Message) { m.SessionID = "" }},
		{name: "empty role", mutate: func(m *message.Message) { m.Role = "" }},
		{name: "empty text", mutate: func(m *message.Message) { m.Text = "" }},
		{name: "zero captured_at", mutate: func(m *message.Message) { m.CapturedAt = time.Time{} }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := valid
			tc.mutate(&msg)

			require.Error(t, msg.Validate())
		})
	}
}
