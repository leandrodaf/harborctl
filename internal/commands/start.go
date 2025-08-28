package commands

import (
	"context"
	"flag"

	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/docker"
)

// startCommand implementa o comando start
type startCommand struct {
	dockerService docker.Service
	output        cli.Output
}

// NewStartCommand cria um novo comando start
func NewStartCommand(dockerService docker.Service, output cli.Output) cli.Command {
	return &startCommand{
		dockerService: dockerService,
		output:        output,
	}
}

func (c *startCommand) Name() string {
	return "start"
}

func (c *startCommand) Description() string {
	return "Inicia serviços previamente parados (docker compose start)"
}

func (c *startCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("start", flag.ExitOnError)

	var outputPath string
	fs.StringVar(&outputPath, "f", ".deploy/compose.generated.yml", "compose file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	c.output.Info("▶️  Iniciando serviços...")
	return c.dockerService.Start(ctx, outputPath)
}
