package natspub

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/AFK068/skilloom/ingest/internal/domain/message"
	"github.com/AFK068/skilloom/ingest/internal/domain/principal"
)

type messagePayload struct {
	EventID    uuid.UUID       `json:"event_id"`
	UserID     uuid.UUID       `json:"user_id"`
	SpaceID    uuid.UUID       `json:"space_id"`
	Harness    string          `json:"harness"`
	SessionID  string          `json:"session_id"`
	TurnID     string          `json:"turn_id,omitempty"`
	Role       string          `json:"role"`
	Text       string          `json:"text"`
	CapturedAt time.Time       `json:"captured_at"`
	Meta       json.RawMessage `json:"meta,omitempty"`
}

func toMessagePayload(p principal.Principal, msg message.Message) messagePayload {
	return messagePayload{
		EventID:    msg.EventID,
		UserID:     p.UserID,
		SpaceID:    p.SpaceID,
		Harness:    msg.Harness,
		SessionID:  msg.SessionID,
		TurnID:     msg.TurnID,
		Role:       msg.Role,
		Text:       msg.Text,
		CapturedAt: msg.CapturedAt,
		Meta:       msg.Meta,
	}
}
