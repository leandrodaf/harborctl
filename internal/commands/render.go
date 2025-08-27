package commands

import (
	"context"
	"flag"

	"github.com/leandrodaf/harborctl/internal/compose"
	"github.com/leandrodaf/harborctl/internal/config"
	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/fs"
)

// renderCommand implementa o comando render
type renderCommand struct {
	configManager  config.Manager
	composeService compose.Service
	filesystem     fs.FileSystem
	output         cli.Output
}

// NewRenderCommand cria um novo comando render
func NewRenderCommand(
	configManager config.Manager,
	composeService compose.Service,
	filesystem fs.FileSystem,
	output cli.Output,
) cli.Command {
	return &renderCommand{
		configManager:  configManager,
		composeService: composeService,
		filesystem:     filesystem,
		output:         output,
	}
}

func (c *renderCommand) Name() string {
	return "render"
}

func (c *renderCommand) Description() string {
	return "Generates docker-compose.yml"
}

func (c *renderCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("render", flag.ExitOnError)

	var stackPath, outputPath string
	var noDozzle, noBeszel bool

	fs.StringVar(&stackPath, "f", "stack.yml", "stack.yml")
	fs.StringVar(&outputPath, "o", ".deploy/compose.generated.yml", "output compose")
	fs.BoolVar(&noDozzle, "no-dozzle", false, "don't include dozzle")
	fs.BoolVar(&noBeszel, "no-beszel", false, "don't include beszel")

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

	// Write file
	if err := c.filesystem.WriteFile(outputPath, data, 0644); err != nil {
		return err
	}

	c.output.Infof("compose generated at %s", outputPath)
	return nil
}
