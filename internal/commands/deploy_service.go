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

	fs.StringVar(&serviceName, "service", "", "microservice name")
	fs.StringVar(&repoURL, "repo", "", "repository URL (optional if already cloned)")
	fs.StringVar(&branch, "branch", "main", "repository branch")
	fs.StringVar(&path, "path", "deploy", "stack.yml path in repository")
	fs.StringVar(&envFile, "env-file", "", "environment variables file")
	fs.StringVar(&secretsFile, "secrets-file", "", "secrets file")
	fs.IntVar(&replicas, "replicas", 0, "number of replicas (override)")
	fs.BoolVar(&dryRun, "dry-run", false, "validate only without deploying")
	fs.BoolVar(&force, "force", false, "force deployment ignoring warnings")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if serviceName == "" {
		return fmt.Errorf("--service is required")
	}

	c.output.Infof("🚀 Deploying microservice: %s", serviceName)

	// Load server base configuration
	baseConfig, err := c.loadBaseConfig(ctx)
	if err != nil {
		c.output.Error("❌ Base configuration not found. Run first:")
		c.output.Error("   harborctl init-server --domain <your-domain> --email <your-email>")
		return err
	}

	// Get microservice code
	serviceDir, err := c.getServiceCode(ctx, serviceName, repoURL, branch, path)
	if err != nil {
		return fmt.Errorf("error getting service code: %w", err)
	}

	// Load microservice configuration
	stackPath := filepath.Join(serviceDir, "stack.yml")
	serviceConfig, err := c.configManager.Load(ctx, stackPath)
	if err != nil {
		return fmt.Errorf("error loading service stack.yml: %w", err)
	}

	// Apply env and secrets overrides
	if err := c.applyRuntimeOverrides(serviceConfig, envFile, secretsFile, replicas); err != nil {
		return fmt.Errorf("error applying overrides: %w", err)
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

	// Deploy microservice
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
		return "", fmt.Errorf("service code not found and --repo not specified")
	}

	// Clone or update repository
	c.output.Infof("📦 Getting code from repository: %s", repoURL)

	if err := c.filesystem.MkdirAll(".services", 0755); err != nil {
		return "", err
	}

	// Verificar se já existe
	if exists := c.filesystem.Exists(serviceDir); exists {
		// Fazer pull
		if err := c.gitClient.Pull(ctx, serviceDir, branch); err != nil {
			return "", fmt.Errorf("error updating repository: %w", err)
		}
		c.output.Info("✅ Repository updated")
	} else {
		// Clonar
		token := getTokenFromEnv("GITHUB_TOKEN") // Tentar obter token
		if err := c.gitClient.Clone(ctx, repoURL, serviceDir, token); err != nil {
			return "", fmt.Errorf("error cloning repository: %w", err)
		}
		c.output.Info("✅ Repository cloned")
	}

	return filepath.Join(serviceDir, path), nil
}

func (c *deployServiceCommand) applyRuntimeOverrides(serviceConfig *config.Stack, envFile, secretsFile string, replicas int) error {
	// Environment variables override
	if envFile != "" {
		c.output.Infof("📝 Applying environment variables from: %s", envFile)
		envVars, err := c.loadEnvFile(envFile)
		if err != nil {
			return fmt.Errorf("erro ao carregar env file: %w", err)
		}

		// Apply vars to all services
		for i := range serviceConfig.Services {
			if serviceConfig.Services[i].Env == nil {
				serviceConfig.Services[i].Env = make(map[string]string)
			}
			for key, value := range envVars {
				serviceConfig.Services[i].Env[key] = value
			}
		}
	}

	// Secrets override
	if secretsFile != "" {
		c.output.Infof("🔐 Applying secrets from: %s", secretsFile)
		// TODO: Implement secrets application
	}

	// Replicas override
	if replicas > 0 {
		c.output.Infof("📈 Setting %d replicas", replicas)
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

		// Only service volumes (base already running)
		Volumes: serviceConfig.Volumes,

		// Only microservice services (infrastructure already running)
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
	c.output.Info("🚢 Deploying microservice...")

	// Generate compose only for microservice
	options := compose.GenerateOptions{
		DisableDozzle: true, // Já está rodando na base
		DisableBeszel: true, // Já está rodando na base
	}

	data, err := c.composeService.Generate(ctx, config, options)
	if err != nil {
		return err
	}

	// Create deploy directory
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
		Prune:  false, // Don't prune to avoid affecting other services
		Detach: true,
	}

	if err := c.dockerService.Deploy(ctx, composePath, deployOptions); err != nil {
		return fmt.Errorf("deployment error: %w", err)
	}

	c.output.Infof("✅ Microservice %s deployed successfully!", serviceName)
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
	// TODO: Implement .env file parser
	// Por enquanto, assumir formato KEY=VALUE

	return envVars, nil
}

func getTokenFromEnv(envVar string) string {
	// TODO: Implement secure environment variables reading
	return ""
}
