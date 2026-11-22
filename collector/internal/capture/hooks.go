package capture

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/AFK068/skilloom/collector/internal/domain/harness"
	"github.com/AFK068/skilloom/collector/internal/domain/message"
	"github.com/AFK068/skilloom/collector/internal/domain/role"
)

type hookInput struct {
	SessionID            string `json:"session_id"`
	TurnID               string `json:"turn_id"`
	TranscriptPath       string `json:"transcript_path"`
	Prompt               string `json:"prompt"`
	LastAssistantMessage string `json:"last_assistant_message"`
}

func readMessage(agent harness.Harness, author role.Role, input io.Reader) (*message.Message, error) {
	var event hookInput
	if err := json.NewDecoder(input).Decode(&event); err != nil {
		return nil, errors.New("could not read hook input")
	}

	var text string
	switch author {
	case role.User:
		text = event.Prompt
	case role.Assistant:
		text = event.LastAssistantMessage
	default:
		return nil, errors.New("unknown message role")
	}
	if event.SessionID == "" || strings.TrimSpace(text) == "" {
		return nil, nil
	}
	if agent == harness.Codex && event.TranscriptPath == "" && isCodexComposerSuggestion(author, text) {
		return nil, nil
	}

	return &message.Message{
		EventID:    uuid.New(),
		Harness:    agent,
		SessionID:  event.SessionID,
		TurnID:     event.TurnID,
		Role:       author,
		Text:       text,
		CapturedAt: time.Now().UTC(),
	}, nil
}

func isCodexComposerSuggestion(author role.Role, text string) bool {
	if author == role.User {
		return strings.HasPrefix(
			text,
			"This is an app-generated request for a suggested next user message. The user did not write this message.",
		)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(text), &fields); err != nil || len(fields) != 1 {
		return false
	}
	var suggestion *string
	return json.Unmarshal(fields["suggestion"], &suggestion) == nil
}
