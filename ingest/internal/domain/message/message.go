package message

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	EventID    uuid.UUID
	Harness    string
	SessionID  string
	TurnID     string
	Role       string
	Text       string
	CapturedAt time.Time
	Meta       json.RawMessage
}
