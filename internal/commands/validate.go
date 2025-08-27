package commands

import (
	"context"
	"flag"

	"github.com/leandrodaf/harborctl/internal/config"
	"github.com/leandrodaf/harborctl/pkg/cli"
)

// validateCommand implementa o comando validate
type validateCommand struct {
	configManager   config.Manager
	secureValidator *config.SecureValidator
	output          cli.Output
}

// NewValidateCommand cria um novo comando validate
func NewValidateCommand(configManager config.Manager, output cli.Output) cli.Command {
	return &validateCommand{
		configManager:   configManager,
		secureValidator: config.NewSecureValidator(),
		output:          output,
	}
}

func (c *validateCommand) Name() string {
	return "validate"
}

func (c *validateCommand) Description() string {
	return "Valida a configuração do stack.yml"
}

func (c *validateCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)

	var stackPath string
	fs.StringVar(&stackPath, "f", "stack.yml", "caminho do stack.yml")

	if err := fs.Parse(args); err != nil {
		return err
	}

	stack, err := c.configManager.Load(ctx, stackPath)
	if err != nil {
		return err
	}

	// Validação padrão
	if err := c.configManager.Validate(ctx, stack); err != nil {
		return err
	}

	// Validação de segurança
	if err := c.secureValidator.ValidateStack(stack); err != nil {
		c.output.Error("Falha na validação de segurança: " + err.Error())
		return err
	}

	c.output.Info("OK: configuração válida e segura.")
	return nil
}
