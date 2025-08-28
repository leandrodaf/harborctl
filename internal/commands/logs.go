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

// logsCommand implementa o comando logs
type logsCommand struct {
	dockerService docker.Service
	output        cli.Output
}

// NewLogsCommand cria um novo comando logs
func NewLogsCommand(dockerService docker.Service, output cli.Output) cli.Command {
	return &logsCommand{
		dockerService: dockerService,
		output:        output,
	}
}

func (c *logsCommand) Name() string {
	return "logs"
}

func (c *logsCommand) Description() string {
	return "Mostra logs dos serviços"
}

func (c *logsCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("logs", flag.ExitOnError)

	var service, composePath string
	var follow bool
	var tail int

	fs.StringVar(&service, "service", "", "nome do serviço")
	fs.StringVar(&composePath, "f", ".deploy/compose.generated.yml", "arquivo compose")
	fs.BoolVar(&follow, "follow", false, "seguir logs em tempo real")
	fs.IntVar(&tail, "tail", 50, "número de linhas para mostrar")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Verificar se o arquivo compose existe
	if !fileExistsLogs(composePath) {
		c.output.Errorf("❌ Arquivo compose não encontrado: %s", composePath)
		c.output.Info("💡 Execute 'harborctl up -f server-base.yml' para criar a infraestrutura")
		return fmt.Errorf("compose file not found: %s", composePath)
	}

	// Preparar comando docker compose logs
	args_cmd := []string{"compose", "-f", composePath, "logs"}

	if follow {
		args_cmd = append(args_cmd, "-f")
	}

	if tail > 0 {
		args_cmd = append(args_cmd, "--tail", fmt.Sprintf("%d", tail))
	}

	if service != "" {
		c.output.Infof("📋 Logs do serviço: %s", service)
		args_cmd = append(args_cmd, service)
	} else {
		c.output.Info("📋 Logs de todos os serviços:")
	}

	// Executar comando
	cmd := exec.CommandContext(ctx, "docker", args_cmd...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func fileExistsLogs(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
