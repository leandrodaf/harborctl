package commands

import (
	"context"
	"flag"

	"github.com/leandrodaf/harborctl/pkg/cli"
)

// docsCommand implementa o comando docs
type docsCommand struct {
	output cli.Output
}

// NewDocsCommand cria um novo comando docs
func NewDocsCommand(output cli.Output) cli.Command {
	return &docsCommand{
		output: output,
	}
}

func (c *docsCommand) Name() string {
	return "docs"
}

func (c *docsCommand) Description() string {
	return "Mostra documentação e exemplos de configuração"
}

func (c *docsCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("docs", flag.ExitOnError)

	var topic string
	fs.StringVar(&topic, "topic", "overview", "tópico (overview, resources, auth, scaling, secrets)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	switch topic {
	case "overview":
		c.showOverview()
	case "resources":
		c.showResources()
	case "auth":
		c.showAuth()
	case "scaling":
		c.showScaling()
	case "secrets":
		c.showSecrets()
	default:
		c.output.Errorf("Tópico desconhecido: %s", topic)
		c.showAvailableTopics()
	}

	return nil
}

func (c *docsCommand) showOverview() {
	c.output.Info("📚 Harbor CTL - Visão Geral")
	c.output.Info("")
	c.output.Info("Comandos principais:")
	c.output.Info("  harborctl init --domain example.com --email admin@example.com")
	c.output.Info("  harborctl validate")
	c.output.Info("  harborctl render")
	c.output.Info("  harborctl up")
	c.output.Info("  harborctl down")
	c.output.Info("")
	c.output.Info("Comandos utilitários:")
	c.output.Info("  harborctl hash-password --generate")
	c.output.Info("  harborctl scale app=3")
	c.output.Info("  harborctl status")
	c.output.Info("")
	c.output.Info("Use --topic para ver documentação específica:")
	c.showAvailableTopics()
}

func (c *docsCommand) showResources() {
	c.output.Info("💾 Configuração de Recursos")
	c.output.Info("")
	c.output.Info("Exemplo de configuração de recursos:")
	c.output.Info(`
  resources:
    memory: "1g"           # Limite de memória
    cpus: "1.0"           # Limite de CPU
    reserve_mem: "512m"   # Reserva de memória
    reserve_cpu: "0.5"    # Reserva de CPU
    shm_size: "128m"      # Tamanho do /dev/shm
    gpus: "1"             # GPU (requer nvidia runtime)
    ulimits:
      nofile:
        soft: 1024
        hard: 2048`)
	c.output.Info("")
	c.output.Info("Formatos válidos:")
	c.output.Info("  Memória: 512m, 1g, 2048M, etc.")
	c.output.Info("  CPU: 0.5, 1.0, 2, etc.")
	c.output.Info("  GPU: 1, 2, all")
}

func (c *docsCommand) showAuth() {
	c.output.Info("🔐 Autenticação Básica")
	c.output.Info("")
	c.output.Info("1. Gerar senha:")
	c.output.Info("   harborctl hash-password --generate")
	c.output.Info("")
	c.output.Info("2. Configurar no stack.yml:")
	c.output.Info(`
  basic_auth:
    enabled: true
    users:
      admin: "$2a$10$..."  # Hash gerado
      user2: "$2a$10$..."
    # OU usar arquivo:
    users_file: "/path/to/htpasswd"`)
	c.output.Info("")
	c.output.Info("O basic auth é aplicado automaticamente no Traefik")
}

func (c *docsCommand) showScaling() {
	c.output.Info("📈 Escalonamento de Serviços")
	c.output.Info("")
	c.output.Info("1. Configurar réplicas no stack.yml:")
	c.output.Info(`
  services:
    - name: app
      replicas: 3  # Número de réplicas`)
	c.output.Info("")
	c.output.Info("2. Escalar dinamicamente:")
	c.output.Info("   harborctl scale app=5")
	c.output.Info("")
	c.output.Info("Load balancing automático:")
	c.output.Info("  - Health checks em /health")
	c.output.Info("  - Sticky sessions desabilitadas")
	c.output.Info("  - Distribuição round-robin")
}

func (c *docsCommand) showSecrets() {
	c.output.Info("🔑 Gerenciamento de Secrets")
	c.output.Info("")
	c.output.Info("Configuração de secrets:")
	c.output.Info(`
  secrets:
    - name: db_password
      file: "./secrets/db_password.txt"
      target: "/run/secrets/db_password"
    - name: api_key
      external: true  # Secret já existe no Docker`)
	c.output.Info("")
	c.output.Info("Uso no container:")
	c.output.Info("  - Arquivos em /run/secrets/")
	c.output.Info("  - Leia via cat /run/secrets/db_password")
	c.output.Info("  - Permissões restritas automaticamente")
}

func (c *docsCommand) showAvailableTopics() {
	c.output.Info("Tópicos disponíveis:")
	c.output.Info("  --topic overview   (visão geral)")
	c.output.Info("  --topic resources  (configuração de recursos)")
	c.output.Info("  --topic auth       (autenticação básica)")
	c.output.Info("  --topic scaling    (escalonamento)")
	c.output.Info("  --topic secrets    (gerenciamento de secrets)")
}
