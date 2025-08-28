package compose

import (
	"context"

	"github.com/leandrodaf/harborctl/internal/config"
)

// NetworkBuilder constrói redes
type NetworkBuilder interface {
	Build(ctx context.Context, networks map[string]config.Network) map[string]map[string]any
}

// VolumeBuilder constrói volumes
type VolumeBuilder interface {
	Build(ctx context.Context, volumes []config.Volume) map[string]map[string]any
}

// ServiceBuilder constrói serviços
type ServiceBuilder interface {
	BuildWithEnvironment(ctx context.Context, service config.Service, domain string, env Environment, project string) map[string]any
}

// TraefikBuilder constrói configuração do Traefik
type TraefikBuilder interface {
	Build(ctx context.Context, stack *config.Stack, env Environment) map[string]any
}

// ObservabilityBuilder constrói serviços de observabilidade
type ObservabilityBuilder interface {
	Build(ctx context.Context, observability config.Observability, domain string, env Environment, options GenerateOptions, project string) map[string]map[string]any
}

// HealthChecker define estratégias de health check
type HealthChecker interface {
	Build(healthConfig config.HealthCheck, port int) map[string]interface{}
}

// DeployStrategy define estratégias de deploy
type DeployStrategy interface {
	Build(deployConfig config.DeployConfig, replicas int) map[string]interface{}
}

// Marshaler serializa compose
type Marshaler interface {
	Marshal(compose *ComposeFile) ([]byte, error)
}
