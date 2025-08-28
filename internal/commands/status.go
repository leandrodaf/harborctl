package commands

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/docker"
)

// statusCommand implements the status command
type statusCommand struct {
	dockerService docker.Service
	output        cli.Output
}

// NewStatusCommand creates a new status command
func NewStatusCommand(dockerService docker.Service, output cli.Output) cli.Command {
	return &statusCommand{
		dockerService: dockerService,
		output:        output,
	}
}

func (c *statusCommand) Name() string {
	return "status"
}

func (c *statusCommand) Description() string {
	return "Show services status"
}

func (c *statusCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("status", flag.ExitOnError)

	var composePath string
	var verbose bool
	fs.StringVar(&composePath, "f", ".deploy/compose.generated.yml", "compose file")
	fs.BoolVar(&verbose, "verbose", false, "show detailed status")

	if err := fs.Parse(args); err != nil {
		return err
	}

	c.output.Info("🔍 Services status:")

	// Check if compose file exists
	if !fileExistsStatus(composePath) {
		c.output.Errorf("❌ Compose file not found: %s", composePath)
		c.output.Info("💡 Run 'harborctl up -f server-base.yml' to create infrastructure")
		return fmt.Errorf("compose file not found: %s", composePath)
	}

	// Executar docker compose ps
	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composePath, "ps", "--format", "table")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		c.output.Error("❌ Error checking container status")
		return fmt.Errorf("failed to get container status: %w", err)
	}

	if verbose {
		c.output.Info("\n🔍 Detailed status:")

		// Show resource statistics
		c.output.Info("\n📊 Resource usage:")
		statsCmd := exec.CommandContext(ctx, "docker", "stats", "--no-stream", "--format",
			"table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}")
		statsCmd.Stdout = os.Stdout
		statsCmd.Run()

		// Show networks
		c.output.Info("\n🌐 Networks:")
		netCmd := exec.CommandContext(ctx, "docker", "network", "ls", "--filter", "name=deploy")
		netCmd.Stdout = os.Stdout
		netCmd.Run()

		// Show volumes
		c.output.Info("\n💾 Volumes:")
		volCmd := exec.CommandContext(ctx, "docker", "volume", "ls", "--filter", "name=deploy")
		volCmd.Stdout = os.Stdout
		volCmd.Run()
	}

	return nil
}

func fileExistsStatus(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
