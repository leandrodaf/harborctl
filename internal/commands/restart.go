package commands

import (
	"context"
	"flag"

	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/docker"
)

// restartCommand implementa o comando restart
type restartCommand struct {
	dockerService docker.Service
	output        cli.Output
}

// NewRestartCommand cria um novo comando restart
func NewRestartCommand(dockerService docker.Service, output cli.Output) cli.Command {
	return &restartCommand{
		dockerService: dockerService,
		output:        output,
	}
}

func (c *restartCommand) Name() string {
	return "restart"
}

func (c *restartCommand) Description() string {
	return "Reinicia todos os serviços (docker compose restart)"
}

func (c *restartCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("restart", flag.ExitOnError)

	var outputPath string
	var timeout int
	fs.StringVar(&outputPath, "f", ".deploy/compose.generated.yml", "compose file")
	fs.IntVar(&timeout, "t", 10, "timeout em segundos para reiniciar os containers")

	if err := fs.Parse(args); err != nil {
		return err
	}

	c.output.Info("🔄 Reiniciando serviços...")
	return c.dockerService.Restart(ctx, outputPath, timeout)
}
