package commands

import (
	"context"
	"flag"
	"fmt"

	"github.com/leandrodaf/harborctl/internal/config"
	"github.com/leandrodaf/harborctl/pkg/cli"
)

// initServerCommand implementa o comando para criar configuração base do servidor
type initServerCommand struct {
	configManager config.Manager
	output        cli.Output
}

// NewInitServerCommand cria um novo comando init-server
func NewInitServerCommand(configManager config.Manager, output cli.Output) cli.Command {
	return &initServerCommand{
		configManager: configManager,
		output:        output,
	}
}

func (c *initServerCommand) Name() string {
	return "init-server"
}

func (c *initServerCommand) Description() string {
	return "Cria configuração base do servidor (infraestrutura, logs, monitoramento)"
}

func (c *initServerCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("init-server", flag.ExitOnError)

	var domain, email, project string
	var replaceExisting bool

	fs.StringVar(&domain, "domain", "", "domínio base (ex.: exemplo.com)")
	fs.StringVar(&email, "email", "", "email para ACME")
	fs.StringVar(&project, "project", "infrastructure", "nome do projeto base")
	fs.BoolVar(&replaceExisting, "replace", false, "substituir configuração existente")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if domain == "" || email == "" {
		c.output.Error("Uso: harborctl init-server --domain <dominio> --email <email>")
		return fmt.Errorf("domain e email são obrigatórios")
	}

	c.output.Info("🏗️  Criando configuração base do servidor...")

	// Verificar se já existe configuração
	if exists, _ := fileExists("server-base.yml"); exists && !replaceExisting {
		c.output.Error("server-base.yml já existe. Use --replace para sobrescrever")
		return fmt.Errorf("configuração base já existe")
	}

	// Criar configuração base do servidor
	baseConfig := c.createBaseServerConfig(domain, email, project)

	// Salvar configuração
	if err := c.configManager.SaveBaseConfig(ctx, "server-base.yml", baseConfig); err != nil {
		return fmt.Errorf("erro ao criar configuração base: %w", err)
	}

	c.output.Info("✅ Configuração base do servidor criada: server-base.yml")
	c.output.Info("📋 Esta configuração inclui:")
	c.output.Info("   • Traefik (reverse proxy + TLS)")
	c.output.Info("   • Dozzle (logs em tempo real)")
	c.output.Info("   • Beszel (monitoramento)")
	c.output.Info("   • Redes e volumes base")
	c.output.Info("")
	c.output.Info("🚀 Deploy da infraestrutura base:")
	c.output.Info("   harborctl up -f server-base.yml")
	c.output.Info("")
	c.output.Info("📦 Para deployar microserviços:")
	c.output.Info("   harborctl deploy-service --service <nome-servico> --repo <url-repo>")

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
			},
			Beszel: config.Beszel{
				Enabled:      true,
				Subdomain:    "monitor",
				DataVolume:   "beszel_data",
				SocketVolume: "beszel_socket",
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
