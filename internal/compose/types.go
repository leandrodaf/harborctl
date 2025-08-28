package compose

import (
	"strings"

	"github.com/leandrodaf/harborctl/internal/config"
)

// ComposeFile representa um arquivo docker-compose
type ComposeFile struct {
	Version  string                    `yaml:"version"`
	Services map[string]map[string]any `yaml:"services"`
	Networks map[string]map[string]any `yaml:"networks"`
	Volumes  map[string]map[string]any `yaml:"volumes"`
	Secrets  map[string]map[string]any `yaml:"secrets,omitempty"`
}

// Environment define o tipo de ambiente
type Environment string

const (
	EnvironmentLocal      Environment = "local"
	EnvironmentProduction Environment = "production"
)

// IsLocalhost verifica se o ambiente é local
func (env Environment) IsLocalhost() bool {
	return env == EnvironmentLocal
}

// GetEnvironmentFromStack obtém o ambiente do stack ou detecta pelo domínio como fallback
func GetEnvironmentFromStack(stack *config.Stack) Environment {
	// Prioriza o valor explícito do environment no stack
	if stack.Environment != "" {
		env := strings.ToLower(stack.Environment)
		if env == "local" || env == "development" || env == "dev" {
			return EnvironmentLocal
		}
		if env == "production" || env == "prod" {
			return EnvironmentProduction
		}
	}

	// Fallback: detecta pelo domínio (compatibilidade com versões antigas)
	domain := stack.Domain
	if domain == "localhost" || domain == "test.local" || domain == "" ||
		strings.HasSuffix(domain, ".local") || strings.HasSuffix(domain, ".localhost") {
		return EnvironmentLocal
	}
	return EnvironmentProduction
}


