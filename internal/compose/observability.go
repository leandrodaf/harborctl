package compose

import (
	"context"
	"fmt"

	"github.com/leandrodaf/harborctl/internal/config"
)

// ObservabilityBuilderImpl implementa ObservabilityBuilder
type ObservabilityBuilderImpl struct{}

// NewObservabilityBuilder cria um novo ObservabilityBuilder
func NewObservabilityBuilder() ObservabilityBuilder {
	return &ObservabilityBuilderImpl{}
}

// Build constrói serviços de observabilidade
func (o *ObservabilityBuilderImpl) Build(ctx context.Context, observability config.Observability, domain string, env Environment, options GenerateOptions, project string) map[string]map[string]any {
	services := make(map[string]map[string]any)

	// Dozzle (log viewer)
	if !options.DisableDozzle && observability.Dozzle.Enabled {
		services["dozzle"] = o.buildDozzle(observability, domain, env, project)
	}

	// Beszel (monitoring)
	if !options.DisableBeszel && observability.Beszel.Enabled {
		hub, agent := o.buildBeszel(observability, domain, env, project)
		services["beszel-hub"] = hub
		services["beszel-agent"] = agent
	}

	return services
}

// buildDozzle constrói o serviço Dozzle
func (o *ObservabilityBuilderImpl) buildDozzle(observability config.Observability, domain string, env Environment, project string) map[string]any {
	var entrypoint string
	var subdomain string
	var labels map[string]string

	if env.IsLocalhost() {
		// Configuração para desenvolvimento local
		entrypoint = "web"
		if domain == "localhost" {
			subdomain = "logs.localhost"
		} else {
			subdomain = fmt.Sprintf("logs.%s", domain)
		}
		labels = map[string]string{
			"traefik.enable":         "true",
			"traefik.docker.network": project + "_traefik",
			"traefik.http.services.dozzle.loadbalancer.server.port": "8080",
			"traefik.http.routers.dozzle.rule":                      fmt.Sprintf("Host(`%s`)", subdomain),
			"traefik.http.routers.dozzle.entrypoints":               entrypoint,
		}
	} else {
		// Configuração para produção
		entrypoint = "websecure"
		subdomain = fmt.Sprintf("logs.%s", domain)
		labels = map[string]string{
			"traefik.enable":         "true",
			"traefik.docker.network": project + "_traefik",
			"traefik.http.services.dozzle.loadbalancer.server.port": "8080",
			"traefik.http.routers.dozzle.rule":                      fmt.Sprintf("Host(`%s`)", subdomain),
			"traefik.http.routers.dozzle.entrypoints":               entrypoint,
			"traefik.http.routers.dozzle.tls":                       "true",
			"traefik.http.routers.dozzle.tls.certresolver":          "letsencrypt",
		}
	}

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
		"networks": []string{"private", "traefik"},
		"restart":  "unless-stopped",
		"labels":   labels,
	}

	return service
}

// buildBeszel constrói os serviços Beszel
func (o *ObservabilityBuilderImpl) buildBeszel(observability config.Observability, domain string, env Environment, project string) (hub, agent map[string]any) {
	var entrypoint string
	var subdomain string
	var hubLabels map[string]string

	if env.IsLocalhost() {
		// Configuração para desenvolvimento local
		entrypoint = "web"
		if domain == "localhost" {
			subdomain = "monitor.localhost"
		} else {
			subdomain = fmt.Sprintf("monitor.%s", domain)
		}
		hubLabels = map[string]string{
			"traefik.enable":         "true",
			"traefik.docker.network": project + "_traefik",
			"traefik.http.services.beszel-hub.loadbalancer.server.port": "8090",
			"traefik.http.routers.beszel-hub.rule":                      fmt.Sprintf("Host(`%s`)", subdomain),
			"traefik.http.routers.beszel-hub.entrypoints":               entrypoint,
		}
	} else {
		// Configuração para produção
		entrypoint = "websecure"
		subdomain = fmt.Sprintf("monitor.%s", domain)
		hubLabels = map[string]string{
			"traefik.enable":         "true",
			"traefik.docker.network": project + "_traefik",
			"traefik.http.services.beszel-hub.loadbalancer.server.port": "8090",
			"traefik.http.routers.beszel-hub.rule":                      fmt.Sprintf("Host(`%s`)", subdomain),
			"traefik.http.routers.beszel-hub.entrypoints":               entrypoint,
			"traefik.http.routers.beszel-hub.tls":                       "true",
			"traefik.http.routers.beszel-hub.tls.certresolver":          "letsencrypt",
		}
	}

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
		"networks": []string{"private", "traefik"},
		"restart":  "unless-stopped",
		"labels":   hubLabels,
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
		"networks": []string{"private", "traefik"},
		"restart":  "unless-stopped",
	}

	return hub, agent
}
