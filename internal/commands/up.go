package commands

import (
	"context"
	"flag"

	"github.com/leandrodaf/harborctl/internal/compose"
	"github.com/leandrodaf/harborctl/internal/config"
	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/docker"
	"github.com/leandrodaf/harborctl/pkg/fs"
)

// upCommand implementa o comando up
type upCommand struct {
	configManager  config.Manager
	composeService compose.Service
	dockerService  docker.Service
	filesystem     fs.FileSystem
	output         cli.Output
}

// NewUpCommand cria um novo comando up
func NewUpCommand(
	configManager config.Manager,
	composeService compose.Service,
	dockerService docker.Service,
	filesystem fs.FileSystem,
	output cli.Output,
) cli.Command {
	return &upCommand{
		configManager:  configManager,
		composeService: composeService,
		dockerService:  dockerService,
		filesystem:     filesystem,
		output:         output,
	}
}

func (c *upCommand) Name() string {
	return "up"
}

func (c *upCommand) Description() string {
	return "Gera compose e executa deploy (render + docker compose up -d --build + prune)"
}

func (c *upCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("up", flag.ExitOnError)

	var stackPath, outputPath string
	var noDozzle, noBeszel bool

	fs.StringVar(&stackPath, "f", "stack.yml", "stack.yml")
	fs.StringVar(&outputPath, "o", ".deploy/compose.generated.yml", "compose de saída")
	fs.BoolVar(&noDozzle, "no-dozzle", false, "não incluir dozzle")
	fs.BoolVar(&noBeszel, "no-beszel", false, "não incluir beszel")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Carregar e validar configuração
	stack, err := c.configManager.Load(ctx, stackPath)
	if err != nil {
		return err
	}

	if err := c.configManager.Validate(ctx, stack); err != nil {
		return err
	}

	// Gerar compose
	options := compose.GenerateOptions{
		DisableDozzle: noDozzle,
		DisableBeszel: noBeszel,
	}

	data, err := c.composeService.Generate(ctx, stack, options)
	if err != nil {
		return err
	}

	// Criar diretório se não existir
	if err := c.filesystem.MkdirAll(".deploy", 0755); err != nil {
		return err
	}

	// Escrever arquivo
	if err := c.filesystem.WriteFile(outputPath, data, 0644); err != nil {
		return err
	}

	c.output.Infof("compose gerado em %s", outputPath)

	// Deploy
	deployOptions := docker.DeployOptions{
		Build:  true,
		Prune:  true,
		Detach: true,
	}

	return c.dockerService.Deploy(ctx, outputPath, deployOptions)
}
