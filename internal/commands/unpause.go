package commands

import (
	"context"
	"flag"

	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/docker"
)

// unpauseCommand implementa o comando unpause
type unpauseCommand struct {
	dockerService docker.Service
	output        cli.Output
}

// NewUnpauseCommand cria um novo comando unpause
func NewUnpauseCommand(dockerService docker.Service, output cli.Output) cli.Command {
	return &unpauseCommand{
		dockerService: dockerService,
		output:        output,
	}
}

func (c *unpauseCommand) Name() string {
	return "unpause"
}

func (c *unpauseCommand) Description() string {
	return "Despausa todos os serviços (docker compose unpause)"
}

func (c *unpauseCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("unpause", flag.ExitOnError)

	var outputPath string
	fs.StringVar(&outputPath, "f", ".deploy/compose.generated.yml", "compose file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	c.output.Info("▶️  Despausando serviços...")
	return c.dockerService.Unpause(ctx, outputPath)
}
