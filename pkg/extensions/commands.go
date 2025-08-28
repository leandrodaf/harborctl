package extensions

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/leandrodaf/harborctl/pkg/cli"
)

// StatusCommand implements a status command
type StatusCommand struct {
	output cli.Output
}

// NewStatusCommand creates a new status command
func NewStatusCommand(output cli.Output) cli.Command {
	return &StatusCommand{
		output: output,
	}
}

func (c *StatusCommand) Name() string {
	return "status"
}

func (c *StatusCommand) Description() string {
	return "Shows services status"
}

func (c *StatusCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("status", flag.ExitOnError)

	var composePath string
	var verbose bool
	fs.StringVar(&composePath, "f", ".deploy/compose.generated.yml", "compose file")
	fs.BoolVar(&verbose, "verbose", false, "show detailed status")

	if err := fs.Parse(args); err != nil {
		return err
	}

	c.output.Info("🔍 Status dos serviços:")

	// Verificar se o arquivo compose existe
	if !fileExists(composePath) {
		c.output.Errorf("❌ Arquivo compose não encontrado: %s", composePath)
		c.output.Info("💡 Execute 'harborctl up' para criar a infraestrutura")
		return fmt.Errorf("compose file not found: %s", composePath)
	}

	// Executar docker compose ps
	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composePath, "ps", "--format", "table")
	output, err := cmd.Output()
	if err != nil {
		c.output.Error("❌ Erro ao verificar status dos containers")
		return fmt.Errorf("failed to get container status: %w", err)
	}

	// Mostrar output do docker compose ps
	statusOutput := strings.TrimSpace(string(output))
	if statusOutput == "" {
		c.output.Info("📭 Nenhum serviço rodando")
		c.output.Info("💡 Execute 'harborctl up -f server-base.yml' para iniciar")
		return nil
	}

	fmt.Println(statusOutput)

	if verbose {
		c.output.Info("\n🔍 Status detalhado:")

		// Mostrar estatísticas de recursos
		statsCmd := exec.CommandContext(ctx, "docker", "stats", "--no-stream", "--format",
			"table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}")
		if statsOutput, err := statsCmd.Output(); err == nil {
			fmt.Println(string(statsOutput))
		}

		// Mostrar redes
		c.output.Info("\n🌐 Redes:")
		netCmd := exec.CommandContext(ctx, "docker", "network", "ls", "--filter", "name=deploy")
		if netOutput, err := netCmd.Output(); err == nil {
			fmt.Println(string(netOutput))
		}

		// Mostrar volumes
		c.output.Info("\n💾 Volumes:")
		volCmd := exec.CommandContext(ctx, "docker", "volume", "ls", "--filter", "name=deploy")
		if volOutput, err := volCmd.Output(); err == nil {
			fmt.Println(string(volOutput))
		}
	}

	// Verificar saúde dos serviços principais
	c.checkServiceHealth()

	return nil
}

func (c *StatusCommand) checkServiceHealth() {
	c.output.Info("\n🏥 Verificando saúde dos serviços:")

	// Verificar Traefik
	if c.pingService("http://localhost/", "Traefik") {
		c.output.Info("  ✅ Traefik: Respondendo")
	} else {
		c.output.Info("  ❌ Traefik: Não responde")
	}

	// Verificar Dozzle
	if c.pingService("http://logs.localhost/", "Dozzle") {
		c.output.Info("  ✅ Dozzle: Respondendo")
	} else {
		c.output.Info("  ❌ Dozzle: Não responde")
	}

	// Verificar Beszel
	if c.pingService("http://monitor.localhost/", "Beszel") {
		c.output.Info("  ✅ Beszel: Respondendo")
	} else {
		c.output.Info("  ❌ Beszel: Não responde")
	}
}

func (c *StatusCommand) pingService(url, name string) bool {
	cmd := exec.Command("curl", "-s", "-o", "/dev/null", "-w", "%{http_code}", "--max-time", "5", url)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	statusCode := strings.TrimSpace(string(output))
	return statusCode == "200" || statusCode == "404" // 404 é OK para Traefik sem rotas
}

func fileExists(path string) bool {
	cmd := exec.Command("test", "-f", path)
	return cmd.Run() == nil
}

// LogsCommand implements a command to view logs
type LogsCommand struct {
	output cli.Output
}

// NewLogsCommand creates a new logs command
func NewLogsCommand(output cli.Output) cli.Command {
	return &LogsCommand{
		output: output,
	}
}

func (c *LogsCommand) Name() string {
	return "logs"
}

func (c *LogsCommand) Description() string {
	return "Shows services logs"
}

func (c *LogsCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("logs", flag.ExitOnError)

	var service, composePath string
	var follow bool
	var tail int

	fs.StringVar(&service, "service", "", "service name")
	fs.StringVar(&composePath, "f", ".deploy/compose.generated.yml", "compose file")
	fs.BoolVar(&follow, "follow", false, "follow logs")
	fs.IntVar(&tail, "tail", 50, "number of lines to show")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Verificar se o arquivo compose existe
	if !fileExists(composePath) {
		c.output.Errorf("❌ Arquivo compose não encontrado: %s", composePath)
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
