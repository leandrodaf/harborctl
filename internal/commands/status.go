package commands

import (
	"context"
	"flag"

	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/docker"
)

// statusCommand implementa o comando status
type statusCommand struct {
	dockerService docker.Service
	output        cli.Output
}

// NewStatusCommand cria um novo comando status
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
	return "Mostra o status dos serviços"
}

func (c *statusCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("status", flag.ExitOnError)

	var composePath string
	fs.StringVar(&composePath, "f", ".deploy/compose.generated.yml", "arquivo compose")

	if err := fs.Parse(args); err != nil {
		return err
	}

	c.output.Info("🔍 Status dos serviços:")
	c.output.Info("")

	// TODO: Implementar via dockerService quando tiver métodos de status
	// Por enquanto, dar uma mensagem informativa
	c.output.Info("Para ver o status atual, execute:")
	c.output.Infof("  docker compose -f %s ps", composePath)
	c.output.Info("")
	c.output.Info("Para logs em tempo real:")
	c.output.Infof("  docker compose -f %s logs -f", composePath)

	return nil
}
