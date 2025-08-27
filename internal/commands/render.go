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
	return "Gera o docker-compose.yml"
}

func (c *renderCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("render", flag.ExitOnError)

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
	return nil
}
