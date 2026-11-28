package message

import (
	"encoding/json"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func (m Message) Validate() error {
	return validation.ValidateStruct(
		&m,
		validation.Field(
			&m.EventID,
			validation.Required,
			validation.NotIn(uuid.Nil.String()).ErrorObject(validation.ErrRequired),
		),
		validation.Field(&m.Harness, validation.Required),
		validation.Field(&m.SessionID, validation.Required),
		validation.Field(&m.Role, validation.Required),
		validation.Field(&m.Text, validation.Required),
		validation.Field(&m.CapturedAt, validation.Required),
	)
}
