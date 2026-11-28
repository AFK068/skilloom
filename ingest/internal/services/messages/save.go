package messages

import (
	"context"
	"errors"
	"fmt"

	"github.com/AFK068/skilloom/ingest/internal/domain/message"
	"github.com/AFK068/skilloom/ingest/internal/domain/principal"
)

var ErrInvalidMessage = errors.New("invalid message")

func (s *Service) SaveMessages(ctx context.Context, p principal.Principal, msgs []message.Message) error {
	if len(msgs) == 0 {
		return fmt.Errorf("%w: empty batch", ErrInvalidMessage)
	}

	for i, msg := range msgs {
		if err := msg.Validate(); err != nil {
			return fmt.Errorf("%w [%d]: %w", ErrInvalidMessage, i, err)
		}
	}

	for _, msg := range msgs {
		if err := s.Publisher.PublishMessage(ctx, p, msg); err != nil {
			return fmt.Errorf("save message %s: %w", msg.EventID, err)
		}
	}

	return nil
}
