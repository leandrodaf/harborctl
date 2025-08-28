package compose

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/leandrodaf/harborctl/internal/config"
)

// ServiceBuilderImpl implementa a construção de serviços com deploy estratégico
type ServiceBuilderImpl struct {
	healthChecker  HealthChecker
	deployStrategy DeployStrategy
}

// NewServiceBuilder cria um novo ServiceBuilder
func NewServiceBuilder(healthChecker HealthChecker, deployStrategy DeployStrategy) ServiceBuilder {
	return &ServiceBuilderImpl{
		healthChecker:  healthChecker,
		deployStrategy: deployStrategy,
	}
}

// BuildWithEnvironment constrói um serviço com ambiente explícito
func (sb *ServiceBuilderImpl) BuildWithEnvironment(ctx context.Context, service config.Service, domain string, env Environment, project string) map[string]any {
	serviceConfig := make(map[string]interface{})

	// Nome do container
	serviceConfig["container_name"] = service.Name

	// Image ou build
	if service.Build != nil {
		buildConfig := map[string]interface{}{
			"context":    service.Build.Context,
			"dockerfile": service.Build.Dockerfile,
		}
		if len(service.Build.Args) > 0 {
			buildConfig["args"] = service.Build.Args
		}
		serviceConfig["build"] = buildConfig
	} else if service.Image != "" {
		serviceConfig["image"] = service.Image
	}

	// Portas
	if service.Expose > 0 {
		serviceConfig["expose"] = []string{strconv.Itoa(service.Expose)}
	}

	// Variáveis de ambiente
	if len(service.Env) > 0 {
		serviceConfig["environment"] = service.Env
	}

	// Arquivos de ambiente
	if len(service.EnvFile) > 0 {
		serviceConfig["env_file"] = service.EnvFile
	}

	// Volumes
	if len(service.Volumes) > 0 {
		volumes := make([]string, len(service.Volumes))
		for i, vol := range service.Volumes {
			volumes[i] = fmt.Sprintf("%s:%s", vol.Source, vol.Target)
		}
		serviceConfig["volumes"] = volumes
	}

	// Secrets
	if len(service.Secrets) > 0 {
		secrets := make([]map[string]interface{}, len(service.Secrets))
		for i, secret := range service.Secrets {
			secretConfig := map[string]interface{}{
				"source": secret.Name,
				"target": secret.Target,
			}
			secrets[i] = secretConfig
		}
		serviceConfig["secrets"] = secrets
	}

	// Resources
	if service.Resources != nil {
		sb.addResourceLimits(serviceConfig, service.Resources)
	}

	// Health check
	if service.HealthCheck != nil && service.HealthCheck.Enabled {
		healthConfig := sb.healthChecker.Build(*service.HealthCheck, service.Expose)
		if healthConfig != nil {
			serviceConfig["healthcheck"] = healthConfig
		}
	}

	// Deploy configuration
	if service.Deploy != nil {
		deployConfig := sb.deployStrategy.Build(*service.Deploy, service.Replicas)
		if deployConfig != nil {
			serviceConfig["deploy"] = deployConfig
		}
	}

	// Labels do Traefik
	if sb.isTraefikEnabled(service) {
		labels := sb.buildTraefikLabelsWithEnvironment(service, domain, env, project)
		serviceConfig["labels"] = labels
	}

	// Networks - serviços ficam em redes privadas por segurança
	networks := sb.buildNetworks(service, env)
	serviceConfig["networks"] = networks

	// Restart policy
	serviceConfig["restart"] = "unless-stopped"

	// Configurações de segurança do container (apenas para produção)
	if !env.IsLocalhost() {
		sb.addSecurityConfig(serviceConfig)
	}

	return serviceConfig
}

// isTraefikEnabled verifica se o Traefik está habilitado (compatibilidade com formato antigo e novo)
func (sb *ServiceBuilderImpl) isTraefikEnabled(service config.Service) bool {
	traefik := service.GetTraefik()
	return traefik != nil && traefik.Enabled
}

// buildTraefikLabelsWithEnvironment constrói as labels do Traefik com ambiente
func (sb *ServiceBuilderImpl) buildTraefikLabelsWithEnvironment(service config.Service, domain string, env Environment, project string) map[string]string {
	labels := make(map[string]string)
	serviceName := service.Name

	// Labels básicas
	labels["traefik.enable"] = "true"
	labels["traefik.docker.network"] = project + "_traefik"

	// Port do serviço
	labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", serviceName)] = strconv.Itoa(service.Expose)

	// Configurações customizadas vs padrão
	var rule string
	var entryPoints []string
	var middlewares []string
	var priority int

	traefik := service.GetTraefik()
	if traefik != nil {
		// Configurações customizadas
		if traefik.Rule != "" {
			rule = traefik.Rule
		} else {
			rule = fmt.Sprintf("Host(`%s.%s`)", service.Subdomain, domain)
		}

		if len(traefik.EntryPoints) > 0 {
			entryPoints = traefik.EntryPoints
		} else {
			// EntryPoints baseados no ambiente
			if env.IsLocalhost() {
				entryPoints = []string{"web"}
			} else {
				entryPoints = []string{"websecure"}
			}
		}

		middlewares = append(middlewares, traefik.Middlewares...)
		priority = traefik.Priority

		// Labels customizados
		for key, value := range traefik.Labels {
			labels[key] = value
		}
	} else {
		// Configuração padrão
		rule = fmt.Sprintf("Host(`%s.%s`)", service.Subdomain, domain)
		if env.IsLocalhost() {
			entryPoints = []string{"web"}
		} else {
			entryPoints = []string{"websecure"}
			middlewares = []string{"security-headers", "rate-limit", "request-size"}
		}
	}

	// Aplicar configurações do router
	routerName := serviceName
	labels[fmt.Sprintf("traefik.http.routers.%s.rule", routerName)] = rule
	labels[fmt.Sprintf("traefik.http.routers.%s.entrypoints", routerName)] = strings.Join(entryPoints, ",")

	if priority > 0 {
		labels[fmt.Sprintf("traefik.http.routers.%s.priority", routerName)] = strconv.Itoa(priority)
	}

	// TLS apenas para produção
	if !env.IsLocalhost() {
		if traefik != nil && traefik.TLS != nil {
			// Configurações TLS customizadas
			labels[fmt.Sprintf("traefik.http.routers.%s.tls", routerName)] = "true"
			if traefik.TLS.CertResolver != "" {
				labels[fmt.Sprintf("traefik.http.routers.%s.tls.certresolver", routerName)] = traefik.TLS.CertResolver
			} else {
				labels[fmt.Sprintf("traefik.http.routers.%s.tls.certresolver", routerName)] = "letsencrypt"
			}
			if traefik.TLS.Options != "" {
				labels[fmt.Sprintf("traefik.http.routers.%s.tls.options", routerName)] = traefik.TLS.Options
			}
		} else {
			// TLS padrão para produção
			labels[fmt.Sprintf("traefik.http.routers.%s.tls", routerName)] = "true"
			labels[fmt.Sprintf("traefik.http.routers.%s.tls.certresolver", routerName)] = "letsencrypt"
		}
	}

	// Middlewares
	if len(middlewares) > 0 {
		// Adicionar middlewares de timeout padrão para produção
		if !env.IsLocalhost() {
			timeoutMiddleware := fmt.Sprintf("%s-timeout", serviceName)
			labels[fmt.Sprintf("traefik.http.middlewares.%s.forwardauth.authresponseheaders", timeoutMiddleware)] = "X-Forwarded-User"
			labels[fmt.Sprintf("traefik.http.middlewares.%s.circuitbreaker.expression", timeoutMiddleware)] = "NetworkErrorRatio() > 0.30"
			middlewares = append(middlewares, timeoutMiddleware)
		}

		labels[fmt.Sprintf("traefik.http.routers.%s.middlewares", routerName)] = strings.Join(middlewares, ",")
	}

	// Configurações de Load Balancer avançadas
	if traefik != nil && traefik.LoadBalancer != nil {
		lb := traefik.LoadBalancer

		if lb.HealthCheck != nil {
			if lb.HealthCheck.Path != "" {
				labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.healthcheck.path", serviceName)] = lb.HealthCheck.Path
			}
			if lb.HealthCheck.Interval != "" {
				labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.healthcheck.interval", serviceName)] = lb.HealthCheck.Interval
			}
			if lb.HealthCheck.Timeout != "" {
				labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.healthcheck.timeout", serviceName)] = lb.HealthCheck.Timeout
			}
		} else {
			// Health check padrão para produção
			if !env.IsLocalhost() {
				labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.healthcheck.timeout", serviceName)] = "10s"
				labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.healthcheck.interval", serviceName)] = "30s"
			}
		}

		if lb.ResponseForwarding != nil && lb.ResponseForwarding.FlushInterval != "" {
			labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.responseforwarding.flushinterval", serviceName)] = lb.ResponseForwarding.FlushInterval
		} else if !env.IsLocalhost() {
			labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.responseforwarding.flushinterval", serviceName)] = "100ms"
		}

		if lb.Sticky != nil && lb.Sticky.Cookie != nil {
			labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.sticky.cookie", serviceName)] = "true"
			if lb.Sticky.Cookie.Name != "" {
				labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.sticky.cookie.name", serviceName)] = lb.Sticky.Cookie.Name
			}
		} else if service.Replicas > 1 {
			labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.sticky.cookie", serviceName)] = "true"
		}
	} else {
		// Configurações padrão de Load Balancer para produção
		if !env.IsLocalhost() {
			labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.healthcheck.timeout", serviceName)] = "10s"
			labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.healthcheck.interval", serviceName)] = "30s"
			labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.responseforwarding.flushinterval", serviceName)] = "100ms"
		}

		if service.Replicas > 1 {
			labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.sticky.cookie", serviceName)] = "true"
		}
	}

	// Basic Auth se habilitado
	if service.BasicAuth != nil && service.BasicAuth.Enabled {
		middlewareName := fmt.Sprintf("%s-auth", serviceName)
		labels[fmt.Sprintf("traefik.http.middlewares.%s.basicauth.users", middlewareName)] = sb.buildBasicAuthUsers(service.BasicAuth)

		// Adicionar auth middleware à lista
		currentMiddlewares := labels[fmt.Sprintf("traefik.http.routers.%s.middlewares", routerName)]
		if currentMiddlewares != "" {
			labels[fmt.Sprintf("traefik.http.routers.%s.middlewares", routerName)] = currentMiddlewares + "," + middlewareName
		} else {
			labels[fmt.Sprintf("traefik.http.routers.%s.middlewares", routerName)] = middlewareName
		}
	}

	return labels
}

// addSecurityConfig adiciona configurações de segurança ao container
func (sb *ServiceBuilderImpl) addSecurityConfig(serviceConfig map[string]interface{}) {
	// Security options
	serviceConfig["security_opt"] = []string{
		"no-new-privileges:true",
	}

	// Remover capacidades desnecessárias
	serviceConfig["cap_drop"] = []string{"ALL"}

	// Adicionar apenas capacidades essenciais se necessário
	serviceConfig["cap_add"] = []string{
		"CHOWN",
		"SETGID",
		"SETUID",
	}

	// Container read-only quando possível (pode ser override por volumes)
	// serviceConfig["read_only"] = true

	// tmpfs para arquivos temporários
	serviceConfig["tmpfs"] = []string{
		"/tmp:rw,noexec,nosuid,size=100m",
		"/var/tmp:rw,noexec,nosuid,size=50m",
	}

	// Limites de memória e file descriptors
	serviceConfig["ulimits"] = map[string]interface{}{
		"nofile": map[string]interface{}{
			"soft": 65536,
			"hard": 65536,
		},
		"nproc": map[string]interface{}{
			"soft": 4096,
			"hard": 4096,
		},
	}

	// User namespace (rodar como non-root quando possível)
	serviceConfig["user"] = "1000:1000"
}

// addResourceLimits adiciona limites de recursos
func (sb *ServiceBuilderImpl) addResourceLimits(serviceConfig map[string]interface{}, resources *config.Resources) {
	deployConfig := make(map[string]interface{})

	if resources.CPUs != "" || resources.Memory != "" {
		resourcesConfig := make(map[string]interface{})

		if resources.CPUs != "" || resources.Memory != "" {
			limits := make(map[string]interface{})
			if resources.CPUs != "" {
				limits["cpus"] = resources.CPUs
			}
			if resources.Memory != "" {
				limits["memory"] = resources.Memory
			}
			resourcesConfig["limits"] = limits
		}

		deployConfig["resources"] = resourcesConfig
	}

	// GPU support
	if resources.GPUs != "" {
		if resources.GPUs == "all" {
			serviceConfig["runtime"] = "nvidia"
			serviceConfig["environment"] = map[string]string{
				"NVIDIA_VISIBLE_DEVICES": "all",
			}
		} else {
			serviceConfig["runtime"] = "nvidia"
			serviceConfig["environment"] = map[string]string{
				"NVIDIA_VISIBLE_DEVICES": resources.GPUs,
			}
		}
	}

	if len(deployConfig) > 0 {
		serviceConfig["deploy"] = deployConfig
	}
}

// buildTraefikLabels constrói as labels do Traefik
func (sb *ServiceBuilderImpl) buildTraefikLabels(service config.Service, domain string) map[string]string {
	labels := make(map[string]string)

	// Detectar ambiente
	env := DetectEnvironment(domain)

	// Labels básicas do Traefik
	labels["traefik.enable"] = "true"
	labels["traefik.docker.network"] = "traefik"

	// Service principal
	serviceName := service.Name
	labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", serviceName)] = strconv.Itoa(service.Expose)

	// Configurações de timeout do load balancer
	labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.healthcheck.timeout", serviceName)] = "10s"
	labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.healthcheck.interval", serviceName)] = "30s"
	labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.responseforwarding.flushinterval", serviceName)] = "100ms"

	// Router principal - configuração baseada no ambiente
	routerName := serviceName
	labels[fmt.Sprintf("traefik.http.routers.%s.rule", routerName)] = fmt.Sprintf("Host(`%s.%s`)", service.Subdomain, domain)

	if env.IsLocalhost() {
		// Ambiente local: HTTP simples
		labels[fmt.Sprintf("traefik.http.routers.%s.entrypoints", routerName)] = "web"
	} else {
		// Ambiente produção: HTTPS com certificados
		labels[fmt.Sprintf("traefik.http.routers.%s.tls", routerName)] = "true"
		labels[fmt.Sprintf("traefik.http.routers.%s.tls.certresolver", routerName)] = "letsencrypt"
		labels[fmt.Sprintf("traefik.http.routers.%s.entrypoints", routerName)] = "websecure"
	}

	// Middleware chain de segurança (apenas produção)
	var middlewares []string
	if !env.IsLocalhost() {
		middlewares = []string{"security-headers", "rate-limit", "request-size"}
	}

	// Basic Auth se habilitado
	if service.BasicAuth != nil && service.BasicAuth.Enabled {
		middlewareName := fmt.Sprintf("%s-auth", serviceName)
		labels[fmt.Sprintf("traefik.http.middlewares.%s.basicauth.users", middlewareName)] = sb.buildBasicAuthUsers(service.BasicAuth)
		middlewares = append(middlewares, middlewareName)
	}

	// Adiciona middleware de timeout específico do serviço (apenas produção)
	if !env.IsLocalhost() {
		timeoutMiddleware := fmt.Sprintf("%s-timeout", serviceName)
		labels[fmt.Sprintf("traefik.http.middlewares.%s.forwardauth.authresponseheaders", timeoutMiddleware)] = "X-Forwarded-User"
		labels[fmt.Sprintf("traefik.http.middlewares.%s.circuitbreaker.expression", timeoutMiddleware)] = "NetworkErrorRatio() > 0.30"
		middlewares = append(middlewares, timeoutMiddleware)
	}

	// Aplica middlewares apenas se houver algum
	if len(middlewares) > 0 {
		labels[fmt.Sprintf("traefik.http.routers.%s.middlewares", routerName)] = strings.Join(middlewares, ",")
	}

	// Load balancing para múltiplas réplicas
	if service.Replicas > 1 {
		labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.sticky.cookie", serviceName)] = "true"
		labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.sticky.cookie.name", serviceName)] = fmt.Sprintf("_%s_server", serviceName)
		labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.sticky.cookie.secure", serviceName)] = "true"
		labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.sticky.cookie.httponly", serviceName)] = "true"
		labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.sticky.cookie.samesite", serviceName)] = "strict"
	}

	return labels
}

// buildBasicAuthUsers constrói a string de usuários para basic auth
func (sb *ServiceBuilderImpl) buildBasicAuthUsers(auth *config.BasicAuth) string {
	var users []string

	// Usuário único (legacy)
	if auth.Username != "" && auth.Password != "" {
		users = append(users, fmt.Sprintf("%s:%s", auth.Username, auth.Password))
	}

	// Múltiplos usuários
	for username, password := range auth.Users {
		users = append(users, fmt.Sprintf("%s:%s", username, password))
	}

	return strings.Join(users, ",")
}

// buildNetworks define as redes apropriadas para o serviço baseado na segurança
func (sb *ServiceBuilderImpl) buildNetworks(service config.Service, env Environment) []string {
	var networks []string

	// Por padrão, serviços ficam em redes privadas (seguro)
	networks = append(networks, "private")

	// Se o serviço tem Traefik habilitado, precisa da rede traefik para roteamento
	if sb.isTraefikEnabled(service) {
		networks = append(networks, "traefik")
	}

	// Verificar se o serviço precisa de acesso à internet
	if service.NetworkAccess != nil {
		if service.NetworkAccess.Internet {
			// Adicionar rede pública para acesso à internet
			networks = append(networks, "public")
		}

		// Adicionar redes customizadas se especificadas
		if len(service.NetworkAccess.Custom) > 0 {
			networks = append(networks, service.NetworkAccess.Custom...)
		}

		// Se forçar apenas interno, remover rede pública
		if service.NetworkAccess.Internal {
			networks = sb.removeNetwork(networks, "public")
		}
	}

	return networks
}

// removeNetwork remove uma rede específica da lista
func (sb *ServiceBuilderImpl) removeNetwork(networks []string, networkToRemove string) []string {
	var result []string
	for _, network := range networks {
		if network != networkToRemove {
			result = append(result, network)
		}
	}
	return result
}
