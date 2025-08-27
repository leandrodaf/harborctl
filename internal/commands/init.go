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
	return "Creates an initial stack.yml"
}

func (c *initCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("init", flag.ExitOnError)

	var domain, email, project string
	var noDozzle, noBeszel bool

	fs.StringVar(&domain, "domain", "", "base domain (ex: example.com)")
	fs.StringVar(&email, "email", "", "email for ACME")
	fs.StringVar(&project, "project", "app", "project name")
	fs.BoolVar(&noDozzle, "no-dozzle", false, "don't include dozzle")
	fs.BoolVar(&noBeszel, "no-beszel", false, "don't include beszel")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if domain == "" || email == "" {
		c.output.Error("Usage: harborctl init --domain <domain> --email <email>")
		return fmt.Errorf("domain and email are required")
	}

	options := config.CreateOptions{
		Domain:   domain,
		Email:    email,
		Project:  project,
		NoDozzle: noDozzle,
		NoBeszel: noBeszel,
	}

	if err := c.configManager.Create(ctx, "stack.yml", options); err != nil {
		if err.Error() == "stack.yml already exists" {
			c.output.Error("stack.yml already exists; won't overwrite")
			os.Exit(1)
		}
		return fmt.Errorf("error creating stack.yml: %w", err)
	}

	c.output.Info("stack.yml created. Edit and add your services in `services:`.")
	return nil
}
