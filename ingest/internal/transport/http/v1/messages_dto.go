package v1

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/AFK068/skilloom/ingest/internal/domain/message"
)

type createMessageRequest struct {
	EventID    uuid.UUID       `json:"event_id"    binding:"required"`
	Harness    string          `json:"harness"     binding:"required"`
	SessionID  string          `json:"session_id"  binding:"required"`
	TurnID     string          `json:"turn_id"`
	Role       string          `json:"role"        binding:"required"`
	Text       string          `json:"text"        binding:"required"`
	CapturedAt time.Time       `json:"captured_at" binding:"required"`
	Meta       json.RawMessage `json:"meta"`
}

func (r createMessageRequest) toDomain() message.Message {
	return message.Message{
		EventID:    r.EventID,
		Harness:    r.Harness,
		SessionID:  r.SessionID,
		TurnID:     r.TurnID,
		Role:       r.Role,
		Text:       r.Text,
		CapturedAt: r.CapturedAt,
		Meta:       r.Meta,
	}
}
