package message

import (
	"time"

	"github.com/google/uuid"

	"github.com/AFK068/skilloom/collector/internal/domain/harness"
	"github.com/AFK068/skilloom/collector/internal/domain/role"
)

type Message struct {
	EventID    uuid.UUID       `json:"event_id"`
	Harness    harness.Harness `json:"harness"`
	SessionID  string          `json:"session_id"`
	TurnID     string          `json:"turn_id,omitempty"`
	Role       role.Role       `json:"role"`
	Text       string          `json:"text"`
	CapturedAt time.Time       `json:"captured_at"`
}
