package natsclient

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

func (c *Client) Publish(ctx context.Context, subject, msgID string, data []byte) error {
	var opts []jetstream.PublishOpt
	if msgID != "" {
		opts = append(opts, jetstream.WithMsgID(msgID))
	}

	if _, err := c.JS.Publish(ctx, subject, data, opts...); err != nil {
		return fmt.Errorf("publish to %q: %w", subject, err)
	}
	return nil
}
