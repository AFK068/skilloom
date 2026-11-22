package capture

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/AFK068/skilloom/collector/internal/agents"
	"github.com/AFK068/skilloom/collector/internal/client"
	"github.com/AFK068/skilloom/collector/internal/config"
	"github.com/AFK068/skilloom/collector/internal/domain/harness"
	"github.com/AFK068/skilloom/collector/internal/domain/role"
)

func Run(agent harness.Harness, author role.Role, input io.Reader) error {
	event, err := agents.ReadHookMessage(agent, author, input)
	if err != nil || event == nil {
		return err
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	cfg, err := config.Load(filepath.Join(dir, "skilloom", "config.json"))
	if err != nil {
		return err
	}

	event.EventID = uuid.New()
	event.CapturedAt = time.Now().UTC()
	return client.Send(cfg.Endpoint, cfg.Token, *event)
}
