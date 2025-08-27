package commands

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/leandrodaf/harborctl/internal/config"
	"github.com/leandrodaf/harborctl/pkg/cli"
)

// initCommand implementa o comando init
type initCommand struct {
	configManager config.Manager
	output        cli.Output
}

// NewInitCommand cria um novo comando init
func NewInitCommand(configManager config.Manager, output cli.Output) cli.Command {
	return &initCommand{
		configManager: configManager,
		output:        output,
	}
}

func (c *initCommand) Name() string {
	return "init"
}

func (c *initCommand) Description() string {
	return "Cria um stack.yml inicial"
}

func (c *initCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("init", flag.ExitOnError)

	var domain, email, project string
	var noDozzle, noBeszel bool

	fs.StringVar(&domain, "domain", "", "domínio base (ex.: exemplo.com)")
	fs.StringVar(&email, "email", "", "email para ACME")
	fs.StringVar(&project, "project", "app", "nome do projeto")
	fs.BoolVar(&noDozzle, "no-dozzle", false, "não incluir dozzle")
	fs.BoolVar(&noBeszel, "no-beszel", false, "não incluir beszel")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if domain == "" || email == "" {
		c.output.Error("Uso: harborctl init --domain <dominio> --email <email>")
		return fmt.Errorf("domain e email são obrigatórios")
	}

	options := config.CreateOptions{
		Domain:   domain,
		Email:    email,
		Project:  project,
		NoDozzle: noDozzle,
		NoBeszel: noBeszel,
	}

	if err := c.configManager.Create(ctx, "stack.yml", options); err != nil {
		if err.Error() == "stack.yml já existe" {
			c.output.Error("stack.yml já existe; não vou sobrescrever")
			os.Exit(1)
		}
		return fmt.Errorf("erro ao criar stack.yml: %w", err)
	}

	c.output.Info("stack.yml criado. Edite e adicione seus serviços em `services:`.")
	return nil
}
