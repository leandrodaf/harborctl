package commands

import (
	"context"
	"flag"

	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/docker"
)

// downCommand implementa o comando down
type downCommand struct {
	dockerService docker.Service
	output        cli.Output
}

// NewDownCommand cria um novo comando down
func NewDownCommand(dockerService docker.Service, output cli.Output) cli.Command {
	return &downCommand{
		dockerService: dockerService,
		output:        output,
	}
}

func (c *downCommand) Name() string {
	return "down"
}

func (c *downCommand) Description() string {
	return "Para e remove os serviços (docker compose down)"
}

func (c *downCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("down", flag.ExitOnError)

	var outputPath string
	fs.StringVar(&outputPath, "f", ".deploy/compose.generated.yml", "compose file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	return c.dockerService.Teardown(ctx, outputPath)
}
