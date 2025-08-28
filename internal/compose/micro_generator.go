package compose

import (
	"context"

	"gopkg.in/yaml.v3"

	"github.com/leandrodaf/harborctl/internal/config"
)

// GeneratorImpl implementa Generator usando micro-interfaces
type GeneratorImpl struct {
	networkBuilder       NetworkBuilder
	volumeBuilder        VolumeBuilder
	serviceBuilder       ServiceBuilder
	traefikBuilder       TraefikBuilder
	observabilityBuilder ObservabilityBuilder
	marshaler            Marshaler
}

// NewMicroGenerator cria um novo generator com micro-interfaces
func NewMicroGenerator() Generator {
	return &GeneratorImpl{
		networkBuilder:       NewNetworkBuilder(),
		volumeBuilder:        NewVolumeBuilder(),
		serviceBuilder:       NewServiceBuilder(NewHealthChecker(), NewDeployStrategy()),
		traefikBuilder:       NewTraefikBuilder(),
		observabilityBuilder: NewObservabilityBuilder(),
		marshaler:            NewMarshaler(),
	}
}

// Generate gera o docker-compose.yml usando micro-interfaces
func (g *GeneratorImpl) Generate(ctx context.Context, stack *config.Stack, options GenerateOptions) ([]byte, error) {
	compose := &ComposeFile{
		Version:  "3.9",
		Services: make(map[string]map[string]any),
		Networks: make(map[string]map[string]any),
		Volumes:  make(map[string]map[string]any),
		Secrets:  make(map[string]map[string]any),
	}

	// Detecta o ambiente baseado na configuração do stack
	env := GetEnvironmentFromStack(stack)

	// Networks
	compose.Networks = g.networkBuilder.Build(ctx, stack.Networks)

	// Adiciona a rede traefik se não existir (necessária para Traefik e observabilidade)
	if _, exists := compose.Networks["traefik"]; !exists {
		compose.Networks["traefik"] = map[string]any{
			"driver": "bridge",
		}
	}

	// Volumes
	compose.Volumes = g.volumeBuilder.Build(ctx, stack.Volumes)

	// Services
	for _, service := range stack.Services {
		serviceConfig := g.serviceBuilder.BuildWithEnvironment(ctx, service, stack.Domain, env, stack.Project)
		compose.Services[service.Name] = serviceConfig
	}

	// Traefik
	traefikConfig := g.traefikBuilder.Build(ctx, stack, env)
	compose.Services["traefik"] = traefikConfig

	// Observability
	if !options.DisableDozzle || !options.DisableBeszel {
		observabilityServices := g.observabilityBuilder.Build(ctx, stack.Observability, stack.Domain, env, options, stack.Project, stack.TLS)
		for name, service := range observabilityServices {
			compose.Services[name] = service
		}
	}

	// Secrets
	g.buildSecrets(compose, stack)

	// Marshal
	return g.marshaler.Marshal(compose)
} // buildSecrets constrói as secrets do compose
func (g *GeneratorImpl) buildSecrets(compose *ComposeFile, stack *config.Stack) {
	secretsMap := make(map[string]bool)

	// Coleta todas as secrets dos serviços
	for _, service := range stack.Services {
		for _, secret := range service.Secrets {
			if !secretsMap[secret.Name] {
				secretConfig := map[string]any{
					"external": secret.External,
				}
				if secret.File != "" {
					secretConfig["file"] = secret.File
				}
				compose.Secrets[secret.Name] = secretConfig
				secretsMap[secret.Name] = true
			}
		}
	}
}

// MarshalerImpl implementa Marshaler
type MarshalerImpl struct{}

// NewMarshaler cria um novo Marshaler
func NewMarshaler() Marshaler {
	return &MarshalerImpl{}
}

// Marshal serializa o compose file para YAML
func (m *MarshalerImpl) Marshal(compose *ComposeFile) ([]byte, error) {
	return yaml.Marshal(compose)
}
