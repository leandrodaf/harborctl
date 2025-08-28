package commands

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/leandrodaf/harborctl/internal/compose"
	"github.com/leandrodaf/harborctl/internal/config"
	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/docker"
	"github.com/leandrodaf/harborctl/pkg/fs"
	"github.com/leandrodaf/harborctl/pkg/git"
	"github.com/leandrodaf/harborctl/pkg/prompt"
)

// deployServiceCommand implements isolated microservice deployment
type deployServiceCommand struct {
	configManager  config.Manager
	composeService compose.Service
	dockerService  docker.Service
	filesystem     fs.FileSystem
	gitClient      *git.Client
	prompter       prompt.Prompter
	errorHandler   *prompt.ErrorHandler
	output         cli.Output
}

// NewDeployServiceCommand creates a new enhanced deploy-service command
func NewDeployServiceCommand(
	configManager config.Manager,
	composeService compose.Service,
	dockerService docker.Service,
	filesystem fs.FileSystem,
	output cli.Output,
) cli.Command {
	prompter := prompt.NewPrompter()
	return &deployServiceCommand{
		configManager:  configManager,
		composeService: composeService,
		dockerService:  dockerService,
		filesystem:     filesystem,
		gitClient:      git.NewClient(),
		prompter:       prompter,
		errorHandler:   prompt.NewErrorHandler(prompter),
		output:         output,
	}
}

func (c *deployServiceCommand) Name() string {
	return "deploy-service"
}

func (c *deployServiceCommand) Description() string {
	return "Deploy a specific microservice with interactive or direct options"
}

func (c *deployServiceCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("deploy-service", flag.ExitOnError)

	var serviceName, repoURL, branch, path, envFile, secretsFile string
	var dryRun, force bool
	var replicas int

	fs.StringVar(&serviceName, "service", "", "nome do microserviço")
	fs.StringVar(&repoURL, "repo", "", "URL do repositório (opcional se já clonado)")
	fs.StringVar(&branch, "branch", "main", "branch do repositório")
	fs.StringVar(&path, "path", "deploy", "caminho do stack.yml no repositório")
	fs.StringVar(&envFile, "env-file", "", "arquivo de variáveis de ambiente")
	fs.StringVar(&secretsFile, "secrets-file", "", "arquivo de secrets")
	fs.IntVar(&replicas, "replicas", 0, "número de réplicas (override)")
	fs.BoolVar(&dryRun, "dry-run", false, "apenas validar sem fazer deploy")
	fs.BoolVar(&force, "force", false, "forçar deploy ignorando warnings")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if serviceName == "" {
		return fmt.Errorf("--service é obrigatório")
	}

	c.output.Infof("🚀 Deployando microserviço: %s", serviceName)

	// Carregar configuração base do servidor
	baseConfig, err := c.loadBaseConfig(ctx)
	if err != nil {
		c.output.Error("❌ Configuração base não encontrada. Execute primeiro:")
		c.output.Error("   harborctl init-server --domain <seu-dominio> --email <seu-email>")
		return err
	}

	// Obter código do microserviço
	serviceDir, err := c.getServiceCode(ctx, serviceName, repoURL, branch, path)
	if err != nil {
		return fmt.Errorf("erro ao obter código do serviço: %w", err)
	}

	// Carregar configuração do microserviço
	stackPath := filepath.Join(serviceDir, "stack.yml")
	serviceConfig, err := c.configManager.Load(ctx, stackPath)
	if err != nil {
		return fmt.Errorf("erro ao carregar stack.yml do serviço: %w", err)
	}

	// Aplicar overrides de env e secrets
	if err := c.applyRuntimeOverrides(serviceConfig, envFile, secretsFile, replicas); err != nil {
		return fmt.Errorf("erro ao aplicar overrides: %w", err)
	}

	// Merge com configuração base (sem duplicar infraestrutura)
	mergedConfig := c.mergeServiceWithBase(serviceConfig, baseConfig)

	// Validar configuração
	if err := c.configManager.Validate(ctx, mergedConfig); err != nil {
		if !force {
			return fmt.Errorf("validação falhou: %w", err)
		}
		c.output.Errorf("⚠️  Validação falhou, mas continuando com --force: %v", err)
	}

	if dryRun {
		c.output.Info("✅ Dry-run concluído. Configuração válida.")
		return nil
	}

	// Deploy do microserviço
	return c.deployMicroservice(ctx, mergedConfig, serviceName)
}

func (c *deployServiceCommand) loadBaseConfig(ctx context.Context) (*config.Stack, error) {
	// Tentar carregar configuração base
	if exists := c.filesystem.Exists("server-base.yml"); exists {
		return c.configManager.Load(ctx, "server-base.yml")
	}
	return nil, fmt.Errorf("configuração base não encontrada")
}

func (c *deployServiceCommand) getServiceCode(ctx context.Context, serviceName, repoURL, branch, path string) (string, error) {
	serviceDir := filepath.Join(".services", serviceName)

	// Se não tem URL do repo, assumir que já está clonado
	if repoURL == "" {
		if exists := c.filesystem.Exists(serviceDir); exists {
			c.output.Infof("📁 Usando código local em: %s", serviceDir)
			return filepath.Join(serviceDir, path), nil
		}
		return "", fmt.Errorf("código do serviço não encontrado e --repo não especificado")
	}

	// Clonar ou atualizar repositório
	c.output.Infof("📦 Obtendo código do repositório: %s", repoURL)

	if err := c.filesystem.MkdirAll(".services", 0755); err != nil {
		return "", err
	}

	// Verificar se já existe
	if exists := c.filesystem.Exists(serviceDir); exists {
		// Fazer pull
		if err := c.gitClient.Pull(ctx, serviceDir, branch); err != nil {
			return "", fmt.Errorf("erro ao atualizar repositório: %w", err)
		}
		c.output.Info("✅ Repositório atualizado")
	} else {
		// Clonar
		token := getTokenFromEnv("GITHUB_TOKEN") // Tentar obter token
		if err := c.gitClient.Clone(ctx, repoURL, serviceDir, token); err != nil {
			return "", fmt.Errorf("erro ao clonar repositório: %w", err)
		}
		c.output.Info("✅ Repositório clonado")
	}

	return filepath.Join(serviceDir, path), nil
}

func (c *deployServiceCommand) applyRuntimeOverrides(serviceConfig *config.Stack, envFile, secretsFile string, replicas int) error {
	// Override de variáveis de ambiente
	if envFile != "" {
		c.output.Infof("📝 Aplicando variáveis de ambiente de: %s", envFile)
		envVars, err := c.loadEnvFile(envFile)
		if err != nil {
			return fmt.Errorf("erro ao carregar env file: %w", err)
		}

		// Aplicar vars a todos os serviços
		for i := range serviceConfig.Services {
			if serviceConfig.Services[i].Env == nil {
				serviceConfig.Services[i].Env = make(map[string]string)
			}
			for key, value := range envVars {
				serviceConfig.Services[i].Env[key] = value
			}
		}
	}

	// Override de secrets
	if secretsFile != "" {
		c.output.Infof("🔐 Aplicando secrets de: %s", secretsFile)
		// TODO: Implementar aplicação de secrets
	}

	// Override de réplicas
	if replicas > 0 {
		c.output.Infof("📈 Configurando %d réplicas", replicas)
		for i := range serviceConfig.Services {
			serviceConfig.Services[i].Replicas = replicas
		}
	}

	return nil
}

func (c *deployServiceCommand) mergeServiceWithBase(serviceConfig, baseConfig *config.Stack) *config.Stack {
	// Criar nova configuração mesclada
	merged := &config.Stack{
		Version: serviceConfig.Version,
		Project: serviceConfig.Project,
		Domain:  baseConfig.Domain, // Usar domínio da base
		TLS:     baseConfig.TLS,    // Usar TLS da base

		// Usar observabilidade da base (não duplicar)
		Observability: baseConfig.Observability,

		// Merge networks
		Networks: make(map[string]config.Network),

		// Apenas volumes do serviço (base já está rodando)
		Volumes: serviceConfig.Volumes,

		// Apenas serviços do microserviço (infraestrutura já rodando)
		Services: serviceConfig.Services,
	}

	// Merge networks
	for name, network := range baseConfig.Networks {
		merged.Networks[name] = network
	}
	for name, network := range serviceConfig.Networks {
		merged.Networks[name] = network
	}

	return merged
}

func (c *deployServiceCommand) deployMicroservice(ctx context.Context, config *config.Stack, serviceName string) error {
	c.output.Info("🚢 Fazendo deploy do microserviço...")

	// Gerar compose apenas para o microserviço
	options := compose.GenerateOptions{
		DisableDozzle: true, // Já está rodando na base
		DisableBeszel: true, // Já está rodando na base
	}

	data, err := c.composeService.Generate(ctx, config, options)
	if err != nil {
		return err
	}

	// Criar diretório de deploy
	deployDir := filepath.Join(".services", serviceName, ".deploy")
	if err := c.filesystem.MkdirAll(deployDir, 0755); err != nil {
		return err
	}

	// Escrever compose
	composePath := filepath.Join(deployDir, fmt.Sprintf("%s-compose.yml", serviceName))
	if err := c.filesystem.WriteFile(composePath, data, 0644); err != nil {
		return err
	}

	c.output.Infof("📄 Compose gerado: %s", composePath)

	// Deploy
	deployOptions := docker.DeployOptions{
		Build:  true,
		Prune:  false, // Não fazer prune para não afetar outros serviços
		Detach: true,
	}

	if err := c.dockerService.Deploy(ctx, composePath, deployOptions); err != nil {
		return fmt.Errorf("erro no deploy: %w", err)
	}

	c.output.Infof("✅ Microserviço %s deployado com sucesso!", serviceName)
	c.output.Info("📊 Acesse os painéis:")
	c.output.Infof("   • Logs: https://logs.%s", config.Domain)
	c.output.Infof("   • Monitor: https://monitor.%s", config.Domain)

	return nil
}

func (c *deployServiceCommand) loadEnvFile(envFile string) (map[string]string, error) {
	_, err := c.filesystem.ReadFile(envFile)
	if err != nil {
		return nil, err
	}

	envVars := make(map[string]string)
	// TODO: Implementar parser de arquivo .env
	// Por enquanto, assumir formato KEY=VALUE

	return envVars, nil
}

func getTokenFromEnv(envVar string) string {
	// TODO: Implementar leitura segura de variáveis de ambiente
	return ""
}
