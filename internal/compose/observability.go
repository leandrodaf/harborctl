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
	// Determinar Docker socket path
	dockerSocket := "/var/run/docker.sock"
	if observability.DockerSocket != "" {
		dockerSocket = observability.DockerSocket
	}

	// Configurar labels do Traefik baseado no ambiente
	subdomain := fmt.Sprintf("monitor.%s", domain)
	hubLabels := o.buildBeszelHubLabels(subdomain, project, env, tls)

	// Beszel Hub - configuração simplificada
	hub = map[string]any{
		"image":          "henrygd/beszel:latest",
		"container_name": "beszel-hub",
		"volumes": []string{
			observability.Beszel.DataVolume + ":/beszel_data",
			observability.Beszel.SocketVolume + ":/beszel_socket",
		},
		"environment": o.buildBeszelHubEnvironment(observability.Beszel, domain, env),
		"networks":    []string{"private", "traefik"},
		"restart":     "unless-stopped",
		"labels":      hubLabels,
	}

	// Beszel Agent - configuração otimizada para socket Unix
	agentEnvironment := o.buildBeszelAgentEnvironment(observability.Beszel, domain, env)
	agentVolumes := []string{
		fmt.Sprintf("%s:/var/run/docker.sock:ro", dockerSocket),
		observability.Beszel.SocketVolume + ":/beszel_socket",
		"./beszel_agent_data:/var/lib/beszel-agent",
		"/etc/os-release:/etc/os-release:ro",
	}

	agent = map[string]any{
		"image":          "henrygd/beszel-agent:latest",
		"container_name": "beszel-agent",
		"environment":    agentEnvironment,
		"volumes":        agentVolumes,
		"networks":       []string{"private"},
		"restart":        "unless-stopped",
		"user":           "0",
		"privileged":     false,
		"security_opt":   []string{"no-new-privileges:true", "apparmor:unconfined"},
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

// buildBeszelHubLabels cria labels do Traefik para o Beszel Hub
func (o *ObservabilityBuilderImpl) buildBeszelHubLabels(subdomain, project string, env Environment, tls config.TLS) map[string]string {
	labels := map[string]string{
		"traefik.enable":         "true",
		"traefik.docker.network": project + "_traefik",
		"traefik.http.services.beszel-hub.loadbalancer.server.port": "8090",
		"traefik.http.routers.beszel-hub.rule":                      fmt.Sprintf("Host(`%s`)", subdomain),
	}

	if env.IsLocalhost() {
		labels["traefik.http.routers.beszel-hub.entrypoints"] = "web"
	} else {
		labels["traefik.http.routers.beszel-hub.entrypoints"] = "web,websecure"
		labels["traefik.http.routers.beszel-hub.tls"] = "true"
		labels["traefik.http.routers.beszel-hub.tls.certresolver"] = tls.Resolver
		labels["traefik.http.routers.beszel-hub.tls.domains[0].main"] = subdomain
	}

	return labels
}

// buildBeszelHubEnvironment cria variáveis de ambiente para o Hub
func (o *ObservabilityBuilderImpl) buildBeszelHubEnvironment(beszel config.Beszel, domain string, env Environment) map[string]string {
	environment := map[string]string{
		"PORT": "8090",
	}

	// APP_URL baseado na configuração ou gerado automaticamente
	var appURL string
	if beszel.AppURL != "" {
		appURL = beszel.AppURL
	} else {
		if env.IsLocalhost() {
			appURL = fmt.Sprintf("http://monitor.%s", domain)
		} else {
			appURL = fmt.Sprintf("https://monitor.%s", domain)
		}
	}
	environment["APP_URL"] = appURL

	// Configurações opcionais
	if beszel.UserCreation {
		environment["USER_CREATION"] = "true"
	}

	return environment
}

// buildBeszelAgentEnvironment cria variáveis de ambiente para o Agent
func (o *ObservabilityBuilderImpl) buildBeszelAgentEnvironment(beszel config.Beszel, domain string, env Environment) map[string]string {
	environment := map[string]string{
		"LISTEN": "/beszel_socket/beszel.sock", // Socket Unix para comunicação local
	}

	// Para mesmo host (nosso caso), usar APENAS socket Unix
	// Ainda precisamos de uma chave mínima para o agent não falhar
	// Usar chave dummy - a autenticação real é feita no painel web do Hub
	if beszel.PublicKey != "" {
		environment["KEY"] = beszel.PublicKey
	} else {
		// Chave dummy para evitar erro de inicialização
		environment["KEY"] = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAdrBaSZ2q0kfjS7RS0WO/WFkEKJjXMF4h3zVO3wg/jN dummy@harborctl"
	}

	return environment
}
