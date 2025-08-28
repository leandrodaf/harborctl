package commands

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/leandrodaf/harborctl/internal/config"
	"github.com/leandrodaf/harborctl/pkg/cli"
)

// securityAuditCommand implementa auditoria de segurança
type securityAuditCommand struct {
	configManager   config.Manager
	secureValidator *config.SecureValidator
	output          cli.Output
}

// NewSecurityAuditCommand cria um novo comando de auditoria
func NewSecurityAuditCommand(configManager config.Manager, output cli.Output) cli.Command {
	return &securityAuditCommand{
		configManager:   configManager,
		secureValidator: config.NewSecureValidator(),
		output:          output,
	}
}

func (c *securityAuditCommand) Name() string {
	return "security-audit"
}

func (c *securityAuditCommand) Description() string {
	return "Executa auditoria de segurança completa nas configurações"
}

func (c *securityAuditCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("security-audit", flag.ExitOnError)

	var stackPath, repoConfigPath string
	var includeRepos bool

	fs.StringVar(&stackPath, "f", "stack.yml", "caminho do stack.yml")
	fs.StringVar(&repoConfigPath, "repos", "repos.yml", "configuração de repositórios")
	fs.BoolVar(&includeRepos, "include-repos", false, "incluir auditoria de repositórios")

	if err := fs.Parse(args); err != nil {
		return err
	}

	c.output.Info("🔒 Iniciando auditoria de segurança...")

	// Auditoria da configuração local
	if err := c.auditLocalConfig(ctx, stackPath); err != nil {
		return err
	}

	// Auditoria de repositórios se solicitado
	if includeRepos {
		if err := c.auditRepositories(ctx, repoConfigPath); err != nil {
			return err
		}
	}

	c.output.Info("✅ Auditoria de segurança concluída!")
	return nil
}

func (c *securityAuditCommand) auditLocalConfig(ctx context.Context, stackPath string) error {
	c.output.Info("🔍 Auditando configuração local...")

	// Carregar configuração
	stack, err := c.configManager.Load(ctx, stackPath)
	if err != nil {
		return fmt.Errorf("erro ao carregar configuração: %w", err)
	}

	// Validação de segurança
	if err := c.secureValidator.ValidateStack(stack); err != nil {
		c.output.Error("❌ Falha na validação de segurança:")
		c.output.Error("   " + err.Error())
		return err
	}

	// Verificações específicas
	securityIssues := c.checkSecurityIssues(stack)

	if len(securityIssues) > 0 {
		c.output.Error("⚠️  Problemas de segurança encontrados:")
		for _, issue := range securityIssues {
			c.output.Error("   • " + issue)
		}
		return fmt.Errorf("encontrados %d problemas de segurança", len(securityIssues))
	}

	c.output.Info("✅ Configuração local aprovada na auditoria")
	return nil
}

func (c *securityAuditCommand) auditRepositories(ctx context.Context, repoConfigPath string) error {
	c.output.Info("🔍 Auditing repository configuration...")

	// Check if repo config file exists
	if _, err := os.Stat(repoConfigPath); os.IsNotExist(err) {
		c.output.Info("⚠️  Repository config file not found, skipping repository audit")
		return nil
	}

	// Basic repository security checks
	issues := c.checkRepositorySecurityIssues(repoConfigPath)

	if len(issues) > 0 {
		c.output.Info("⚠️  Repository security issues found:")
		for _, issue := range issues {
			c.output.Info(fmt.Sprintf("   • %s", issue))
		}
	} else {
		c.output.Info("✅ Repository configuration passed security audit")
	}

	return nil
}

func (c *securityAuditCommand) checkRepositorySecurityIssues(repoConfigPath string) []string {
	var issues []string

	// Basic validation checks
	content, err := os.ReadFile(repoConfigPath)
	if err != nil {
		issues = append(issues, fmt.Sprintf("Cannot read repo config: %v", err))
		return issues
	}

	configStr := string(content)

	// Check for insecure protocols
	if strings.Contains(configStr, "http://") {
		issues = append(issues, "Insecure HTTP URLs found in repository config")
	}

	// Check for hardcoded tokens (basic patterns)
	suspiciousPatterns := []string{
		"token=", "password=", "secret=", "key=",
		"ghp_", "gho_", "ghu_", "ghs_", // GitHub tokens
	}

	for _, pattern := range suspiciousPatterns {
		if strings.Contains(strings.ToLower(configStr), pattern) {
			issues = append(issues, fmt.Sprintf("Potential hardcoded credential detected: %s", pattern))
		}
	}

	return issues
}

func (c *securityAuditCommand) checkSecurityIssues(stack *config.Stack) []string {
	var issues []string

	// Verificar TLS
	if stack.TLS.Mode == "disabled" {
		issues = append(issues, "TLS está desabilitado - recomenda-se usar ACME")
	}

	// Verificar serviços
	for _, service := range stack.Services {
		// Verificar exposição pública sem autenticação
		traefik := service.GetTraefik()
		if traefik != nil && traefik.Enabled && service.BasicAuth != nil && !service.BasicAuth.Enabled {
			issues = append(issues, fmt.Sprintf("Serviço '%s' está exposto publicamente sem autenticação", service.Name))
		}

		// Verificar recursos ilimitados
		if service.Resources == nil {
			issues = append(issues, fmt.Sprintf("Serviço '%s' não tem limites de recursos definidos", service.Name))
		}

		// Verificar secrets não externos
		for _, secret := range service.Secrets {
			if !secret.External && secret.File != "" {
				issues = append(issues, fmt.Sprintf("Secret '%s' do serviço '%s' usa arquivo local - considere usar secrets externos", secret.Name, service.Name))
			}
		}

		// Verificar imagens sem tag específica
		if service.Image != "" && !containsTag(service.Image) {
			issues = append(issues, fmt.Sprintf("Serviço '%s' usa imagem sem tag específica - use tags fixas em produção", service.Name))
		}
	}

	return issues
}

func containsTag(image string) bool {
	// Verificação simples se a imagem tem uma tag específica
	return len(image) > 0 && (image[len(image)-1] != ':' && image != "latest")
}
