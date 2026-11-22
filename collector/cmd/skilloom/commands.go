package main

import (
	"fmt"
	"os"

	"github.com/AFK068/skilloom/collector/internal/capture"
	"github.com/AFK068/skilloom/collector/internal/domain/harness"
	"github.com/AFK068/skilloom/collector/internal/domain/role"
)

type CLI struct {
	Capture captureCmd `cmd:"" hidden:""`
}

type captureCmd struct {
	Agent harness.Harness `arg:"" enum:"codex,claude"`
	Role  role.Role       `arg:"" enum:"user,assistant"`
}

func (c *captureCmd) Run() error {
	if err := capture.Run(c.Agent, c.Role, os.Stdin); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return nil
}
