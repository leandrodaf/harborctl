package commands

import (
	"context"
	"flag"

	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/docker"
)

// stopCommand implementa o comando stop
type stopCommand struct {
	dockerService docker.Service
	output        cli.Output
}

// NewStopCommand cria um novo comando stop
func NewStopCommand(dockerService docker.Service, output cli.Output) cli.Command {
	return &stopCommand{
		dockerService: dockerService,
		output:        output,
	}
}

func (c *stopCommand) Name() string {
	return "stop"
}

func (c *stopCommand) Description() string {
	return "Para os serviços sem remover containers (docker compose stop)"
}

func (c *stopCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("stop", flag.ExitOnError)

	var outputPath string
	var timeout int
	fs.StringVar(&outputPath, "f", ".deploy/compose.generated.yml", "compose file")
	fs.IntVar(&timeout, "t", 10, "timeout em segundos para parar os containers")

	if err := fs.Parse(args); err != nil {
		return err
	}

	c.output.Info("⏹️  Parando serviços...")
	return c.dockerService.Stop(ctx, outputPath, timeout)
}
