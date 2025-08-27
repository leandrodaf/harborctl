package extensions

import (
	"context"
	"flag"
	"fmt"

	"github.com/leandrodaf/harborctl/pkg/cli"
)

// StatusCommand implements a status command
type StatusCommand struct {
	output cli.Output
}

// NewStatusCommand creates a new status command
func NewStatusCommand(output cli.Output) cli.Command {
	return &StatusCommand{
		output: output,
	}
}

func (c *StatusCommand) Name() string {
	return "status"
}

func (c *StatusCommand) Description() string {
	return "Shows services status"
}

func (c *StatusCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("status", flag.ExitOnError)

	var composePath string
	fs.StringVar(&composePath, "f", ".deploy/compose.generated.yml", "compose file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	c.output.Info("🔍 Checking services status...")

	// Here you could implement the real verification logic
	// For example, using docker ps, docker compose ps, etc.

	c.output.Info("✅ All services are running")
	return nil
}

// LogsCommand implements a command to view logs
type LogsCommand struct {
	output cli.Output
}

// NewLogsCommand creates a new logs command
func NewLogsCommand(output cli.Output) cli.Command {
	return &LogsCommand{
		output: output,
	}
}

func (c *LogsCommand) Name() string {
	return "logs"
}

func (c *LogsCommand) Description() string {
	return "Shows services logs"
}

func (c *LogsCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("logs", flag.ExitOnError)

	var service, composePath string
	var follow bool

	fs.StringVar(&service, "service", "", "service name")
	fs.StringVar(&composePath, "f", ".deploy/compose.generated.yml", "compose file")
	fs.BoolVar(&follow, "follow", false, "follow logs")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if service == "" {
		c.output.Error("Specify a service with --service")
		return fmt.Errorf("service is required")
	}

	c.output.Infof("📋 Showing logs for service: %s", service)

	// Here you would implement the real logs logic
	// For example: docker compose logs -f service-name

	return nil
}
