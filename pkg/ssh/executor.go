package ssh

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type executor struct{}

// NewExecutor creates a new SSH executor
func NewExecutor() Executor {
	return &executor{}
}

func (e *executor) Execute(ctx context.Context, config Config, command string) error {
	args := []string{
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "LogLevel=ERROR",
		"-p", fmt.Sprintf("%d", config.Port),
	}

	if config.KeyFile != "" {
		args = append(args, "-i", config.KeyFile)
	}

	target := fmt.Sprintf("%s@%s", config.User, config.Host)
	args = append(args, target, command)

	cmd := exec.CommandContext(ctx, "ssh", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
