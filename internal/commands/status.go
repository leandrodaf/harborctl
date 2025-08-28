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
	var verbose bool
	fs.StringVar(&composePath, "f", ".deploy/compose.generated.yml", "arquivo compose")
	fs.BoolVar(&verbose, "verbose", false, "mostrar status detalhado")

	if err := fs.Parse(args); err != nil {
		return err
	}

	c.output.Info("🔍 Status dos serviços:")

	// Verificar se o arquivo compose existe
	if !fileExistsStatus(composePath) {
		c.output.Errorf("❌ Arquivo compose não encontrado: %s", composePath)
		c.output.Info("💡 Execute 'harborctl up -f server-base.yml' para criar a infraestrutura")
		return fmt.Errorf("compose file not found: %s", composePath)
	}

	// Executar docker compose ps
	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composePath, "ps", "--format", "table")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		c.output.Error("❌ Erro ao verificar status dos containers")
		return fmt.Errorf("failed to get container status: %w", err)
	}

	if verbose {
		c.output.Info("\n🔍 Status detalhado:")

		// Mostrar estatísticas de recursos
		c.output.Info("\n📊 Uso de recursos:")
		statsCmd := exec.CommandContext(ctx, "docker", "stats", "--no-stream", "--format",
			"table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}")
		statsCmd.Stdout = os.Stdout
		statsCmd.Run()

		// Mostrar redes
		c.output.Info("\n🌐 Redes:")
		netCmd := exec.CommandContext(ctx, "docker", "network", "ls", "--filter", "name=deploy")
		netCmd.Stdout = os.Stdout
		netCmd.Run()

		// Mostrar volumes
		c.output.Info("\n💾 Volumes:")
		volCmd := exec.CommandContext(ctx, "docker", "volume", "ls", "--filter", "name=deploy")
		volCmd.Stdout = os.Stdout
		volCmd.Run()
	}

	return nil
}

func fileExistsStatus(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
