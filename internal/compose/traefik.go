package compose

import (
	"context"
	"fmt"
	"strings"

	"github.com/leandrodaf/harborctl/internal/config"
)

// traefikBuilder implementa TraefikBuilder
type traefikBuilder struct{}

func NewTraefikBuilder() TraefikBuilder {
	return &traefikBuilder{}
}

func (b *traefikBuilder) Build(ctx context.Context, stack *config.Stack, env Environment) map[string]any {
	var args []string
	var ports []string
	var labels map[string]string
	var volumes []string
	var environment map[string]string

	// Configurações base vs customizadas
	if stack.Traefik != nil {
		// Usar configurações customizadas
		return b.buildCustomTraefik(stack, env)
	}

	// Configurações padrão baseadas no ambiente
	if env.IsLocalhost() {
		args = []string{
			"--providers.docker=true",
			"--providers.docker.exposedbydefault=false",
			"--entrypoints.web.address=:80",
			"--entrypoints.websecure.address=:443",
			"--api.dashboard=true",
			"--api.insecure=true",
			"--log.level=INFO",
			"--providers.docker.network=" + stack.Project + "_traefik",
			"--global.checknewversion=false",
			"--global.sendanonymoususage=false",
		}
		ports = []string{"80:80", "443:443", "8080:8080"}
		labels = map[string]string{
			"traefik.enable": "false",
		}
		volumes = []string{"/var/run/docker.sock:/var/run/docker.sock:ro"}
	} else {
		args = []string{
			"--providers.docker=true",
			"--providers.docker.exposedbydefault=false",
			"--entrypoints.web.address=:80",
			"--entrypoints.websecure.address=:443",
			"--entrypoints.websecure.http.tls=true",
			"--entrypoints.web.http.redirections.entrypoint.to=websecure",
			"--entrypoints.web.http.redirections.entrypoint.scheme=https",
			"--entrypoints.web.http.redirections.entrypoint.permanent=true",
			"--providers.docker.network=" + stack.Project + "_traefik",
			"--global.checknewversion=false",
			"--global.sendanonymoususage=false",
			"--entrypoints.websecure.transport.respondingtimeouts.readtimeout=60s",
			"--entrypoints.websecure.transport.respondingtimeouts.writetimeout=60s",
			"--entrypoints.websecure.transport.respondingtimeouts.idletimeout=180s",
		}
		ports = []string{"80:80", "443:443"}
		labels = map[string]string{
			"traefik.enable": "false",
		}
		volumes = []string{"/var/run/docker.sock:/var/run/docker.sock:ro"}
	}

	// Adiciona configurações ACME apenas se estiver em modo ACME
	if stack.TLS.Mode == "acme" && !env.IsLocalhost() {
		res := stack.TLS.Resolver
		args = append(args,
			fmt.Sprintf("--certificatesresolvers.%s.acme.email=%s", res, stack.TLS.Email),
			fmt.Sprintf("--certificatesresolvers.%s.acme.storage=/letsencrypt/acme.json", res),
		)

		if stack.TLS.DNS != nil && stack.TLS.DNS.Provider != "" {
			// Usar DNS challenge quando configurado
			args = append(args,
				fmt.Sprintf("--certificatesresolvers.%s.acme.dnschallenge=true", res),
				fmt.Sprintf("--certificatesresolvers.%s.acme.dnschallenge.provider=%s", res, stack.TLS.DNS.Provider),
			)

			// Adicionar variáveis de ambiente do DNS provider se fornecidas
			if len(stack.TLS.DNS.Env) > 0 {
				if environment == nil {
					environment = make(map[string]string)
				}
				for _, envVar := range stack.TLS.DNS.Env {
					// Assume formato KEY=VALUE
					parts := strings.SplitN(envVar, "=", 2)
					if len(parts) == 2 {
						environment[parts[0]] = parts[1]
					}
				}
			}
		} else {
			// Usar HTTP challenge como fallback
			args = append(args,
				fmt.Sprintf("--certificatesresolvers.%s.acme.httpchallenge=true", res),
				fmt.Sprintf("--certificatesresolvers.%s.acme.httpchallenge.entrypoint=web", res),
			)
		}
	}

	config := map[string]any{
		"image":    "traefik:v3.5",
		"command":  args,
		"ports":    ports,
		"labels":   labels,
		"networks": []string{"public", "private", "traefik"},
		"restart":  "always",
		"volumes":  volumes,
	}

	if len(environment) > 0 {
		config["environment"] = environment
	}

	// Adiciona volume para ACME apenas em produção
	if stack.TLS.Mode == "acme" && !env.IsLocalhost() {
		volumes = append(volumes, "traefik_acme:/letsencrypt")
		config["volumes"] = volumes
	}

	// Adiciona configurações de segurança apenas em produção
	if !env.IsLocalhost() {
		config["security_opt"] = []string{"no-new-privileges:true"}
		config["read_only"] = true
		config["tmpfs"] = []string{"/tmp:rw,noexec,nosuid,size=100m"}
		config["deploy"] = map[string]any{
			"resources": map[string]any{
				"limits": map[string]string{
					"cpus":   "1.0",
					"memory": "512M",
				},
				"reservations": map[string]string{
					"cpus":   "0.25",
					"memory": "128M",
				},
			},
		}
	}

	return config
}

// buildCustomTraefik constrói configuração customizada do Traefik
func (b *traefikBuilder) buildCustomTraefik(stack *config.Stack, env Environment) map[string]any {
	traefikConfig := stack.Traefik

	// Imagem
	image := "traefik:v3.5"
	if traefikConfig.Image != "" {
		image = traefikConfig.Image
	}

	// Commands
	var commands []string
	if len(traefikConfig.Commands) > 0 {
		commands = traefikConfig.Commands
	} else {
		// Commands padrão baseados no ambiente
		if env.IsLocalhost() {
			commands = []string{
				"--providers.docker=true",
				"--providers.docker.exposedbydefault=false",
				"--entrypoints.web.address=:80",
				"--entrypoints.websecure.address=:443",
				"--api.dashboard=true",
				"--api.insecure=true",
				"--log.level=INFO",
				"--providers.docker.network=" + stack.Project + "_traefik",
				"--global.checknewversion=false",
				"--global.sendanonymoususage=false",
			}
		} else {
			commands = []string{
				"--providers.docker=true",
				"--providers.docker.exposedbydefault=false",
				"--entrypoints.web.address=:80",
				"--entrypoints.websecure.address=:443",
				"--entrypoints.websecure.http.tls=true",
				"--entrypoints.web.http.redirections.entrypoint.to=websecure",
				"--entrypoints.web.http.redirections.entrypoint.scheme=https",
				"--entrypoints.web.http.redirections.entrypoint.permanent=true",
				"--providers.docker.network=" + stack.Project + "_traefik",
				"--global.checknewversion=false",
				"--global.sendanonymoususage=false",
				"--entrypoints.websecure.transport.respondingtimeouts.readtimeout=60s",
				"--entrypoints.websecure.transport.respondingtimeouts.writetimeout=60s",
				"--entrypoints.websecure.transport.respondingtimeouts.idletimeout=180s",
			}
		}
	}

	// Adicionar configurações de entry points customizados
	for name, ep := range traefikConfig.EntryPoints {
		commands = append(commands, fmt.Sprintf("--entrypoints.%s.address=%s", name, ep.Address))
		if ep.AsDefault {
			commands = append(commands, fmt.Sprintf("--entrypoints.%s.asDefault=true", name))
		}
	}

	// Adicionar configurações de providers customizados
	for name, provider := range traefikConfig.Providers {
		if provider.Docker != nil {
			commands = append(commands, fmt.Sprintf("--providers.docker.%s=%v", name, provider.Docker))
		}
		if provider.File != nil {
			commands = append(commands, fmt.Sprintf("--providers.file.%s=%v", name, provider.File))
		}
	}

	// Adicionar configurações de middlewares customizados
	for name, middleware := range traefikConfig.Middlewares {
		commands = append(commands, b.buildMiddlewareCommands(name, middleware)...)
	}

	// Adicionar configurações de plugins
	for name, plugin := range traefikConfig.Plugins {
		commands = append(commands, fmt.Sprintf("--experimental.plugins.%s.modulename=%s", name, plugin.ModuleName))
		if plugin.Version != "" {
			commands = append(commands, fmt.Sprintf("--experimental.plugins.%s.version=%s", name, plugin.Version))
		}
	}

	// API
	if traefikConfig.API != nil {
		if traefikConfig.API.Dashboard {
			commands = append(commands, "--api.dashboard=true")
		}
		if traefikConfig.API.Insecure {
			commands = append(commands, "--api.insecure=true")
		}
		if traefikConfig.API.Debug {
			commands = append(commands, "--api.debug=true")
		}
	}

	// Log
	if traefikConfig.Log != nil {
		if traefikConfig.Log.Level != "" {
			commands = append(commands, fmt.Sprintf("--log.level=%s", traefikConfig.Log.Level))
		}
		if traefikConfig.Log.Format != "" {
			commands = append(commands, fmt.Sprintf("--log.format=%s", traefikConfig.Log.Format))
		}
		if traefikConfig.Log.FilePath != "" {
			commands = append(commands, fmt.Sprintf("--log.filepath=%s", traefikConfig.Log.FilePath))
		}
	}

	// Access Log
	if traefikConfig.AccessLog != nil {
		commands = append(commands, "--accesslog=true")
		if traefikConfig.AccessLog.FilePath != "" {
			commands = append(commands, fmt.Sprintf("--accesslog.filepath=%s", traefikConfig.AccessLog.FilePath))
		}
		if traefikConfig.AccessLog.Format != "" {
			commands = append(commands, fmt.Sprintf("--accesslog.format=%s", traefikConfig.AccessLog.Format))
		}
	}

	// Metrics
	if traefikConfig.Metrics != nil {
		if traefikConfig.Metrics.Prometheus != nil {
			commands = append(commands, "--metrics.prometheus=true")
			if traefikConfig.Metrics.Prometheus.AddEntryPointsLabels {
				commands = append(commands, "--metrics.prometheus.addentrypointslabels=true")
			}
			if traefikConfig.Metrics.Prometheus.AddServicesLabels {
				commands = append(commands, "--metrics.prometheus.addserviceslabels=true")
			}
		}
	}

	// Ports
	ports := []string{"80:80", "443:443", "8080:8080"}
	if len(traefikConfig.Ports) > 0 {
		ports = traefikConfig.Ports
	}

	// Labels
	labels := map[string]string{"traefik.enable": "false"}
	if len(traefikConfig.Labels) > 0 {
		for k, v := range traefikConfig.Labels {
			labels[k] = v
		}
	}

	// Volumes
	volumes := []string{"/var/run/docker.sock:/var/run/docker.sock:ro"}
	if len(traefikConfig.Volumes) > 0 {
		volumes = append(volumes, traefikConfig.Volumes...)
	}

	config := map[string]any{
		"image":    image,
		"command":  commands,
		"ports":    ports,
		"labels":   labels,
		"networks": []string{"public", "private", "traefik"},
		"restart":  "always",
		"volumes":  volumes,
	}

	// Environment variables
	if len(traefikConfig.Environment) > 0 {
		config["environment"] = traefikConfig.Environment
	}

	// Adiciona configurações de segurança apenas em produção
	if !env.IsLocalhost() {
		config["security_opt"] = []string{"no-new-privileges:true"}
		config["read_only"] = true
		config["tmpfs"] = []string{"/tmp:rw,noexec,nosuid,size=100m"}
	}

	return config
}

// buildMiddlewareCommands constrói comandos para middlewares customizados
func (b *traefikBuilder) buildMiddlewareCommands(name string, middleware config.TraefikMiddleware) []string {
	var commands []string

	if middleware.AddPrefix != nil {
		commands = append(commands, fmt.Sprintf("--http.middlewares.%s.addprefix.prefix=%s", name, middleware.AddPrefix.Prefix))
	}

	if middleware.StripPrefix != nil {
		for _, prefix := range middleware.StripPrefix.Prefixes {
			commands = append(commands, fmt.Sprintf("--http.middlewares.%s.stripprefix.prefixes=%s", name, prefix))
		}
		if middleware.StripPrefix.ForceSlash {
			commands = append(commands, fmt.Sprintf("--http.middlewares.%s.stripprefix.forceslash=true", name))
		}
	}

	if middleware.ReplacePathRegex != nil {
		commands = append(commands, fmt.Sprintf("--http.middlewares.%s.replacepathregex.regex=%s", name, middleware.ReplacePathRegex.Regex))
		commands = append(commands, fmt.Sprintf("--http.middlewares.%s.replacepathregex.replacement=%s", name, middleware.ReplacePathRegex.Replacement))
	}

	if middleware.RateLimit != nil {
		if middleware.RateLimit.Average > 0 {
			commands = append(commands, fmt.Sprintf("--http.middlewares.%s.ratelimit.average=%d", name, middleware.RateLimit.Average))
		}
		if middleware.RateLimit.Period != "" {
			commands = append(commands, fmt.Sprintf("--http.middlewares.%s.ratelimit.period=%s", name, middleware.RateLimit.Period))
		}
		if middleware.RateLimit.Burst > 0 {
			commands = append(commands, fmt.Sprintf("--http.middlewares.%s.ratelimit.burst=%d", name, middleware.RateLimit.Burst))
		}
	}

	return commands
}
