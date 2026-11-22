package capture

import (
	"fmt"
	"io"

	"github.com/AFK068/skilloom/collector/internal/client"
	"github.com/AFK068/skilloom/collector/internal/config"
	"github.com/AFK068/skilloom/collector/internal/domain/harness"
	"github.com/AFK068/skilloom/collector/internal/domain/role"
)

func Run(agent harness.Harness, author role.Role, input io.Reader) error {
	event, err := readMessage(agent, author, input)
	if err != nil {
		return fmt.Errorf("read message: %w", err)
	}
	if event == nil {
		return nil
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return client.Send(cfg.Endpoint, cfg.Token, *event)
}
