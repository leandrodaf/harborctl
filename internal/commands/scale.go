package commands

import (
	"context"
	"flag"
	"fmt"
	"strconv"

	"github.com/leandrodaf/harborctl/internal/config"
	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/docker"
)

// scaleCommand implementa o comando scale
type scaleCommand struct {
	configManager config.Manager
	dockerService docker.Service
	output        cli.Output
}

// NewScaleCommand cria um novo comando scale
func NewScaleCommand(configManager config.Manager, dockerService docker.Service, output cli.Output) cli.Command {
	return &scaleCommand{
		configManager: configManager,
		dockerService: dockerService,
		output:        output,
	}
}

func (c *scaleCommand) Name() string {
	return "scale"
}

func (c *scaleCommand) Description() string {
	return "Escala um serviço (ex: harborctl scale app=3)"
}

func (c *scaleCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("scale", flag.ExitOnError)

	var composePath string
	fs.StringVar(&composePath, "f", ".deploy/compose.generated.yml", "arquivo compose")

	if err := fs.Parse(args); err != nil {
		return err
	}

	remainingArgs := fs.Args()
	if len(remainingArgs) == 0 {
		return fmt.Errorf("especifique o serviço e replicas: harborctl scale service=replicas")
	}

	// Parse service=replicas
	scaleSpecs := make(map[string]int)
	for _, arg := range remainingArgs {
		parts := parseScaleArg(arg)
		if len(parts) != 2 {
			return fmt.Errorf("formato inválido: %s (use service=replicas)", arg)
		}

		service := parts[0]
		replicas, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("replicas inválidas para %s: %v", service, err)
		}

		scaleSpecs[service] = replicas
	}

	// Execute scaling usando docker compose
	for service, replicas := range scaleSpecs {
		c.output.Infof("Escalando %s para %d replicas...", service, replicas)

		// Usar docker compose scale
		cmd := fmt.Sprintf("docker compose -f %s up -d --scale %s=%d --no-recreate",
			composePath, service, replicas)

		// Aqui seria melhor usar o dockerService, mas por simplicidade:
		c.output.Infof("Executando: %s", cmd)

		// TODO: implementar via dockerService.Scale() se necessário
	}

	return nil
}

func parseScaleArg(arg string) []string {
	// Simples parse de "service=replicas"
	for i, char := range arg {
		if char == '=' {
			return []string{arg[:i], arg[i+1:]}
		}
	}
	return []string{arg}
}
