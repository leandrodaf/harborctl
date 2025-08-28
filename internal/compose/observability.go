package compose

import (
	"context"
	"fmt"
	"strings"

	"github.com/leandrodaf/harborctl/internal/config"
)

type ObservabilityBuilderImpl struct{}

func NewObservabilityBuilder() ObservabilityBuilder {
	return &ObservabilityBuilderImpl{}
}

func (o *ObservabilityBuilderImpl) Build(ctx context.Context, observability config.Observability, domain string, env Environment, options GenerateOptions, project string, tls config.TLS) map[string]map[string]any {
	services := make(map[string]map[string]any)

	// Dozzle (log viewer)
	if !options.DisableDozzle && observability.Dozzle.Enabled {
		services["dozzle"] = o.buildDozzle(observability, domain, env, project, tls)
	}

	// Beszel (monitoring)
	if !options.DisableBeszel && observability.Beszel.Enabled {
		hub, agent := o.buildBeszel(observability, domain, env, project, tls)
		services["beszel-hub"] = hub
		services["beszel-agent"] = agent
	}

	return services
}

func (o *ObservabilityBuilderImpl) buildDozzle(observability config.Observability, domain string, env Environment, project string, tls config.TLS) map[string]any {
	var entrypoint string
	var subdomain string
	var labels map[string]string

	// Determinar Docker socket path
	dockerSocket := "/var/run/docker.sock"
	if observability.DockerSocket != "" {
		dockerSocket = observability.DockerSocket
	}

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
			"traefik.http.routers.dozzle.tls.certresolver":          tls.Resolver,
		}
	}

	// Adicionar basic auth se configurado
	if observability.Dozzle.BasicAuth != nil && observability.Dozzle.BasicAuth.Enabled {
		middlewareName := "dozzle-auth"
		labels["traefik.http.routers.dozzle.middlewares"] = middlewareName
		labels[fmt.Sprintf("traefik.http.middlewares.%s.basicauth.users", middlewareName)] = o.buildBasicAuthUsers(observability.Dozzle.BasicAuth)
	}

	service := map[string]any{
		"image":          "amir20/dozzle:latest",
		"container_name": "dozzle",
		"volumes": []string{
			fmt.Sprintf("%s:/var/run/docker.sock:ro", dockerSocket),
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

func (o *ObservabilityBuilderImpl) buildBeszel(observability config.Observability, domain string, env Environment, project string, tls config.TLS) (hub, agent map[string]any) {
	var entrypoint string
	var subdomain string
	var hubLabels map[string]string

	// Determinar Docker socket path
	dockerSocket := "/var/run/docker.sock"
	if observability.DockerSocket != "" {
		dockerSocket = observability.DockerSocket
	}

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
		// Configuração para produção baseada na documentação oficial
		entrypoint = "websecure"
		subdomain = fmt.Sprintf("monitor.%s", domain)
		hubLabels = map[string]string{
			"traefik.enable":         "true",
			"traefik.docker.network": project + "_traefik",
			"traefik.http.services.beszel-hub.loadbalancer.server.port": "8090",
			"traefik.http.routers.beszel-hub.rule":                      fmt.Sprintf("Host(`%s`)", subdomain),
			"traefik.http.routers.beszel-hub.entrypoints":               "web,websecure", // Suportar ambos HTTP e HTTPS
			"traefik.http.routers.beszel-hub.tls":                       "true",
			"traefik.http.routers.beszel-hub.tls.certresolver":          tls.Resolver,
			"traefik.http.routers.beszel-hub.tls.domains[0].main":       subdomain,
		}
	}

	// NOTA: Beszel tem sistema de autenticação próprio que conflita com basic auth do Traefik
	// Desabilitando basic auth do Traefik para evitar loops de autenticação
	// O Beszel gerenciará sua própria autenticação internamente

	// Commented out to fix authentication loop issue:
	// if observability.Beszel.BasicAuth != nil && observability.Beszel.BasicAuth.Enabled {
	//     middlewareName := "beszel-auth"
	//     hubLabels["traefik.http.routers.beszel-hub.middlewares"] = middlewareName
	//     hubLabels[fmt.Sprintf("traefik.http.middlewares.%s.basicauth.users", middlewareName)] = o.buildBasicAuthUsers(observability.Beszel.BasicAuth)
	// }

	// Beszel Hub
	hub = map[string]any{
		"image":          "henrygd/beszel:latest",
		"container_name": "beszel-hub",
		"volumes": []string{
			observability.Beszel.DataVolume + ":/beszel_data",
			observability.Beszel.SocketVolume + ":/beszel_socket",
		},
		"environment": buildBeszelHubEnvironment(observability.Beszel, domain, env),
		"networks":    []string{"private", "traefik"},
		"restart":     "unless-stopped",
		"labels":      hubLabels,
	}

	// Beszel Agent - configuração seguindo documentação oficial
	agentEnvironment := map[string]string{
		"LISTEN": "/beszel_socket/beszel.sock", // Usar socket Unix para comunicação local
	}

	// Configurar HUB_URL - usar nome do container Docker
	agentEnvironment["HUB_URL"] = "http://beszel-hub:8090"

	// Configurar token se fornecido
	if observability.Beszel.Token != "" {
		agentEnvironment["TOKEN"] = observability.Beszel.Token
	} else {
		agentEnvironment["TOKEN"] = "CONFIGURE_TOKEN_IN_BESZEL_CONFIG"
	}

	// Configurar chave pública para autenticação
	if observability.Beszel.HubKey != "" {
		agentEnvironment["KEY"] = observability.Beszel.HubKey
	} else if observability.Beszel.HubKeyFile != "" {
		agentEnvironment["KEY_FILE"] = observability.Beszel.HubKeyFile
	} else {
		// Aviso: sem chave configurada, o agent falhará
		agentEnvironment["KEY"] = "CONFIGURE_HUB_KEY_IN_BESZEL_CONFIG"
	}

	// Configurar HUB_URL personalizada se fornecida
	if observability.Beszel.HubURL != "" {
		agentEnvironment["HUB_URL"] = observability.Beszel.HubURL
	}

	agentVolumes := []string{
		fmt.Sprintf("%s:/var/run/docker.sock:ro", dockerSocket),
		observability.Beszel.SocketVolume + ":/beszel_socket",
		"./beszel_agent_data:/var/lib/beszel-agent",
		"/etc/os-release:/etc/os-release:ro", // Para informações do OS
	}

	agent = map[string]any{
		"image":          "henrygd/beszel-agent:latest",
		"container_name": "beszel-agent",
		"environment":    agentEnvironment,
		"volumes":        agentVolumes,
		"networks":       []string{"private"}, // Usar rede Docker normal para conectividade
		"restart":        "unless-stopped",
		"user":           "0", // Executar como root para acessar Docker socket
		"privileged":     false,
		"security_opt":   []string{"no-new-privileges:true"},
		// Removido "pid": "host" para evitar conflito com AppArmor
		// O agent ainda consegue coletar a maioria das estatísticas via /proc e /sys
	}

	return hub, agent
}

func (o *ObservabilityBuilderImpl) buildBasicAuthUsers(auth *config.BasicAuth) string {
	var users []string

	// Usuário único (legacy)
	if auth.Username != "" && auth.Password != "" {
		// Escapar o caractere $ para Docker Compose duplicando-o
		escapedPassword := strings.ReplaceAll(auth.Password, "$", "$$")
		users = append(users, fmt.Sprintf("%s:%s", auth.Username, escapedPassword))
	}

	// Múltiplos usuários
	for username, password := range auth.Users {
		// Escapar o caractere $ para Docker Compose duplicando-o
		escapedPassword := strings.ReplaceAll(password, "$", "$$")
		users = append(users, fmt.Sprintf("%s:%s", username, escapedPassword))
	}

	return strings.Join(users, ",")
}

func buildBeszelHubEnvironment(beszel config.Beszel, domain string, env Environment) map[string]string {
	environment := map[string]string{
		"PORT": "8090",
	}

	// APP_URL baseado na documentação oficial
	var appURL string
	if beszel.AppURL != "" {
		appURL = beszel.AppURL
	} else {
		// Gerar APP_URL automaticamente
		if env.IsLocalhost() {
			if domain == "localhost" {
				appURL = "http://monitor.localhost"
			} else {
				appURL = fmt.Sprintf("http://monitor.%s", domain)
			}
		} else {
			appURL = fmt.Sprintf("https://monitor.%s", domain)
		}
	}
	environment["APP_URL"] = appURL

	// Configurações avançadas de autenticação
	if beszel.DisablePasswordAuth {
		environment["DISABLE_PASSWORD_AUTH"] = "true"
	}

	if beszel.UserCreation {
		environment["USER_CREATION"] = "true"
	}

	return environment
}
