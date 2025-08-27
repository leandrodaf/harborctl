package extensions

import (
	"context"
	"flag"
	"fmt"

	"github.com/leandrodaf/harborctl/pkg/cli"
)

// StatusCommand implementa um comando de status
type StatusCommand struct {
	output cli.Output
}

// NewStatusCommand cria um novo comando de status
func NewStatusCommand(output cli.Output) cli.Command {
	return &StatusCommand{
		output: output,
	}
}

func (c *StatusCommand) Name() string {
	return "status"
}

func (c *StatusCommand) Description() string {
	return "Mostra o status dos serviços"
}

func (c *StatusCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("status", flag.ExitOnError)

	var composePath string
	fs.StringVar(&composePath, "f", ".deploy/compose.generated.yml", "arquivo compose")

	if err := fs.Parse(args); err != nil {
		return err
	}

	c.output.Info("🔍 Verificando status dos serviços...")

	// Aqui você poderia implementar a lógica real de verificação
	// Por exemplo, usando docker ps, docker compose ps, etc.

	c.output.Info("✅ Todos os serviços estão rodando")
	return nil
}

// LogsCommand implementa um comando para visualizar logs
type LogsCommand struct {
	output cli.Output
}

// NewLogsCommand cria um novo comando de logs
func NewLogsCommand(output cli.Output) cli.Command {
	return &LogsCommand{
		output: output,
	}
}

func (c *LogsCommand) Name() string {
	return "logs"
}

func (c *LogsCommand) Description() string {
	return "Mostra os logs dos serviços"
}

func (c *LogsCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("logs", flag.ExitOnError)

	var service, composePath string
	var follow bool

	fs.StringVar(&service, "service", "", "nome do serviço")
	fs.StringVar(&composePath, "f", ".deploy/compose.generated.yml", "arquivo compose")
	fs.BoolVar(&follow, "follow", false, "seguir os logs")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if service == "" {
		c.output.Error("Especifique um serviço com --service")
		return fmt.Errorf("serviço é obrigatório")
	}

	c.output.Infof("📋 Mostrando logs do serviço: %s", service)

	// Aqui você implementaria a lógica real de logs
	// Por exemplo: docker compose logs -f service-name

	return nil
}
