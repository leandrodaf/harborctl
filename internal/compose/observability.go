package compose

import (
	"context"

	"github.com/leandrodaf/harborctl/internal/config"
)

// ObservabilityBuilderImpl implementa ObservabilityBuilder
type ObservabilityBuilderImpl struct{}

// NewObservabilityBuilder cria um novo ObservabilityBuilder
func NewObservabilityBuilder() ObservabilityBuilder {
	return &ObservabilityBuilderImpl{}
}

// Build constrói serviços de observabilidade
func (o *ObservabilityBuilderImpl) Build(ctx context.Context, observability config.Observability, options GenerateOptions) map[string]map[string]any {
	services := make(map[string]map[string]any)

	// Dozzle (log viewer)
	if !options.DisableDozzle && observability.Dozzle.Enabled {
		services["dozzle"] = o.buildDozzle(observability)
	}

	// Beszel (monitoring)
	if !options.DisableBeszel && observability.Beszel.Enabled {
		hub, agent := o.buildBeszel(observability)
		services["beszel-hub"] = hub
		services["beszel-agent"] = agent
	}

	return services
}

// buildDozzle constrói o serviço Dozzle
func (o *ObservabilityBuilderImpl) buildDozzle(observability config.Observability) map[string]any {
	service := map[string]any{
		"image":          "amir20/dozzle:latest",
		"container_name": "dozzle",
		"volumes": []string{
			"/var/run/docker.sock:/var/run/docker.sock:ro",
			observability.Dozzle.DataVolume + ":/data",
		},
		"environment": map[string]string{
			"DOZZLE_LEVEL":    "info",
			"DOZZLE_TAILSIZE": "300",
		},
		"networks": []string{"traefik"},
		"restart":  "unless-stopped",
		"labels": map[string]string{
			"traefik.enable":         "true",
			"traefik.docker.network": "traefik",
			"traefik.http.services.dozzle.loadbalancer.server.port": "8080",
			"traefik.http.routers.dozzle.rule":                      "Host(`logs.` + os.Getenv(\"DOMAIN\"))",
			"traefik.http.routers.dozzle.tls":                       "true",
			"traefik.http.routers.dozzle.tls.certresolver":          "letsencrypt",
			"traefik.http.routers.dozzle.entrypoints":               "websecure",
		},
	}

	return service
}

// buildBeszel constrói os serviços Beszel
func (o *ObservabilityBuilderImpl) buildBeszel(observability config.Observability) (hub, agent map[string]any) {
	// Beszel Hub
	hub = map[string]any{
		"image":          "henrygd/beszel:latest",
		"container_name": "beszel-hub",
		"volumes": []string{
			observability.Beszel.DataVolume + ":/beszel_data",
		},
		"environment": map[string]string{
			"PORT": "8090",
		},
		"networks": []string{"traefik"},
		"restart":  "unless-stopped",
		"labels": map[string]string{
			"traefik.enable":         "true",
			"traefik.docker.network": "traefik",
			"traefik.http.services.beszel-hub.loadbalancer.server.port": "8090",
			"traefik.http.routers.beszel-hub.rule":                      "Host(`monitor.` + os.Getenv(\"DOMAIN\"))",
			"traefik.http.routers.beszel-hub.tls":                       "true",
			"traefik.http.routers.beszel-hub.tls.certresolver":          "letsencrypt",
			"traefik.http.routers.beszel-hub.entrypoints":               "websecure",
		},
	}

	// Beszel Agent
	agent = map[string]any{
		"image":          "henrygd/beszel-agent:latest",
		"container_name": "beszel-agent",
		"environment": map[string]string{
			"PORT": "45876",
		},
		"volumes": []string{
			"/var/run/docker.sock:/var/run/docker.sock:ro",
		},
		"networks": []string{"traefik"},
		"restart":  "unless-stopped",
	}

	return hub, agent
}
