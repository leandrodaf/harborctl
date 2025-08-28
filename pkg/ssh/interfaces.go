package ssh

import "context"

// Config represents SSH connection configuration
type Config struct {
	Host    string
	User    string
	KeyFile string
	Port    int
}

// Executor executes commands via SSH
type Executor interface {
	Execute(ctx context.Context, config Config, command string) error
}

// CommandBuilder builds remote commands
type CommandBuilder interface {
	BuildLogsCommand(composePath, service string, follow bool, tail int) string
	BuildControlCommand(action, composePath, service string, verbose bool) string
}
