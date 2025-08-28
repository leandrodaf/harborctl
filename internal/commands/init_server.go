package commands

import (
	"context"
	"flag"
	"fmt"

	"github.com/leandrodaf/harborctl/internal/config"
	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/prompt"
)

// initServerCommand implements server base configuration creation
type initServerCommand struct {
	configManager config.Manager
	prompter      prompt.Prompter
	errorHandler  *prompt.ErrorHandler
	output        cli.Output
}

// NewInitServerCommand creates a new init-server command
func NewInitServerCommand(configManager config.Manager, output cli.Output) cli.Command {
	prompter := prompt.NewPrompter()
	return &initServerCommand{
		configManager: configManager,
		prompter:      prompter,
		errorHandler:  prompt.NewErrorHandler(prompter),
		output:        output,
	}
}

func (c *initServerCommand) Name() string {
	return "init-server"
}

func (c *initServerCommand) Description() string {
	return "Create server base configuration (infrastructure, logs, monitoring)"
}

func (c *initServerCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("init-server", flag.ExitOnError)

	var domain, email, project string
	var replaceExisting bool

	fs.StringVar(&domain, "domain", "", "base domain (ex: example.com)")
	fs.StringVar(&email, "email", "", "email for ACME certificates")
	fs.StringVar(&project, "project", "infrastructure", "project name")
	fs.BoolVar(&replaceExisting, "replace", false, "replace existing configuration")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if domain == "" || email == "" {
		c.output.Error("Usage: harborctl init-server --domain <domain> --email <email>")
		return fmt.Errorf("domain and email are required")
	}

	// Validate inputs
	if err := prompt.ValidateDomain(domain); err != nil {
		return fmt.Errorf("invalid domain: %w", err)
	}
	
	if err := prompt.ValidateEmail(email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	c.output.Info("🏗️ Creating server base configuration...")

	// Check if configuration already exists
	if exists, _ := fileExists("server-base.yml"); exists && !replaceExisting {
		c.output.Error("server-base.yml already exists. Use --replace to overwrite")
		return fmt.Errorf("base configuration already exists")
	}

	// Create server base configuration
	baseConfig := c.createBaseServerConfig(domain, email, project)

	// Save configuration
	if err := c.configManager.SaveBaseConfig(ctx, "server-base.yml", baseConfig); err != nil {
		return fmt.Errorf("error creating base configuration: %w", err)
	}

	c.output.Info("✅ Server base configuration created: server-base.yml")
	c.output.Info("📋 This configuration includes:")
	c.output.Info("   • Traefik (reverse proxy + TLS)")
	c.output.Info("   • Dozzle (real-time logs)")
	c.output.Info("   • Beszel (monitoring)")
	c.output.Info("   • Base networks and volumes")
	c.output.Info("")
	c.output.Info("🚀 Deploy base infrastructure:")
	c.output.Info("   harborctl up -f server-base.yml")
	c.output.Info("")
	c.output.Info("📦 To deploy microservices:")
	c.output.Info("   harborctl deploy-service --service <service-name> --repo <repo-url>")

	return nil
}

func (c *initServerCommand) createBaseServerConfig(domain, email, project string) *config.Stack {
	return &config.Stack{
		Version: 1,
		Project: project,
		Domain:  domain,
		TLS: config.TLS{
			Mode:     "acme",
			Email:    email,
			Resolver: "le",
		},
		Observability: config.Observability{
			Dozzle: config.Dozzle{
				Enabled:    true,
				Subdomain:  "logs",
				DataVolume: "dozzle_data",
				BasicAuth: &config.BasicAuth{
					Enabled:  true,
					Username: "admin",
					Password: "$2a$10$rO7V2T9JhgHGGkYJlVzZJu.HKVzZqO5qJ5MF5KsGzOzVzSjI2tG6W", // "admin"
				},
			},
			Beszel: config.Beszel{
				Enabled:      true,
				Subdomain:    "monitor",
				DataVolume:   "beszel_data",
				SocketVolume: "beszel_socket",
				BasicAuth: &config.BasicAuth{
					Enabled:  true,
					Username: "admin",
					Password: "$2a$10$rO7V2T9JhgHGGkYJlVzZJu.HKVzZqO5qJ5MF5KsGzOzVzSjI2tG6W", // "admin"
				},
			},
		},
		Networks: map[string]config.Network{
			"private": {Internal: true},
			"public":  {Internal: false},
		},
		Volumes: []config.Volume{
			{Name: "traefik_acme"},
			{Name: "dozzle_data"},
			{Name: "beszel_data"},
			{Name: "beszel_socket"},
		},
		Services: []config.Service{}, // Sem serviços específicos - apenas infraestrutura
	}
}

func fileExists(path string) (bool, error) {
	// TODO: Implementar verificação de arquivo real
	return false, nil
}
