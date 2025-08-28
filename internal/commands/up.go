package commands

import (
	"context"
	"flag"

	"github.com/leandrodaf/harborctl/internal/compose"
	"github.com/leandrodaf/harborctl/internal/config"
	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/docker"
	"github.com/leandrodaf/harborctl/pkg/fs"
)

// upCommand implements the up command
type upCommand struct {
	configManager  config.Manager
	composeService compose.Service
	dockerService  docker.Service
	filesystem     fs.FileSystem
	output         cli.Output
}

// NewUpCommand creates a new up command
func NewUpCommand(
	configManager config.Manager,
	composeService compose.Service,
	dockerService docker.Service,
	filesystem fs.FileSystem,
	output cli.Output,
) cli.Command {
	return &upCommand{
		configManager:  configManager,
		composeService: composeService,
		dockerService:  dockerService,
		filesystem:     filesystem,
		output:         output,
	}
}

func (c *upCommand) Name() string {
	return "up"
}

func (c *upCommand) Description() string {
	return "Generate compose and deploy (render + docker compose up -d --build + prune)"
}

func (c *upCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("up", flag.ExitOnError)

	var stackPath, outputPath string
	var noDozzle, noBeszel bool

	fs.StringVar(&stackPath, "f", "stack.yml", "stack.yml")
	fs.StringVar(&outputPath, "o", ".deploy/compose.generated.yml", "output compose file")
	fs.BoolVar(&noDozzle, "no-dozzle", false, "don't include dozzle")
	fs.BoolVar(&noBeszel, "no-beszel", false, "don't include beszel")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Load and validate configuration
	stack, err := c.configManager.Load(ctx, stackPath)
	if err != nil {
		return err
	}

	if err := c.configManager.Validate(ctx, stack); err != nil {
		return err
	}

	// Generate compose
	options := compose.GenerateOptions{
		DisableDozzle: noDozzle,
		DisableBeszel: noBeszel,
	}

	data, err := c.composeService.Generate(ctx, stack, options)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	if err := c.filesystem.MkdirAll(".deploy", 0755); err != nil {
		return err
	}

	// Write file
	if err := c.filesystem.WriteFile(outputPath, data, 0644); err != nil {
		return err
	}

	c.output.Infof("compose generated at %s", outputPath)

	// Deploy
	deployOptions := docker.DeployOptions{
		Build:  true,
		Prune:  true,
		Detach: true,
	}

	return c.dockerService.Deploy(ctx, outputPath, deployOptions)
}
