package commands

import (
	"context"
	"flag"

	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/docker"
)

// pauseCommand implementa o comando pause
type pauseCommand struct {
	dockerService docker.Service
	output        cli.Output
}

// NewPauseCommand cria um novo comando pause
func NewPauseCommand(dockerService docker.Service, output cli.Output) cli.Command {
	return &pauseCommand{
		dockerService: dockerService,
		output:        output,
	}
}

func (c *pauseCommand) Name() string {
	return "pause"
}

func (c *pauseCommand) Description() string {
	return "Pausa todos os serviços (docker compose pause)"
}

func (c *pauseCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("pause", flag.ExitOnError)

	var outputPath string
	fs.StringVar(&outputPath, "f", ".deploy/compose.generated.yml", "compose file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	c.output.Info("⏸️  Pausando serviços...")
	return c.dockerService.Pause(ctx, outputPath)
}
