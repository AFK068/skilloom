package natspub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AFK068/skilloom/ingest/internal/domain/message"
	"github.com/AFK068/skilloom/ingest/internal/domain/principal"
)

const messagesSubject = "ingest.messages"

func (p *Publisher) PublishMessage(ctx context.Context, pr principal.Principal, msg message.Message) error {
	data, err := json.Marshal(toMessagePayload(pr, msg))
	if err != nil {
		return fmt.Errorf("marshal message %s: %w", msg.EventID, err)
	}

	if err := p.Client.Publish(ctx, messagesSubject, msg.EventID.String(), data); err != nil {
		return fmt.Errorf("publish message %s: %w", msg.EventID, err)
	}

	return nil
}
