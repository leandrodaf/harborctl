package config

// Stack representa a configuração completa
type Stack struct {
	Version       int                `yaml:"version"`
	Project       string             `yaml:"project"`
	Domain        string             `yaml:"domain"`
	Environment   string             `yaml:"environment"` // local | production
	TLS           TLS                `yaml:"tls"`
	Traefik       *TraefikConfig     `yaml:"traefik,omitempty"`
	Observability Observability      `yaml:"observability"`
	Networks      map[string]Network `yaml:"networks"`
	Volumes       []Volume           `yaml:"volumes"`
	Services      []Service          `yaml:"services"`
}

// TLS configura SSL/TLS
type TLS struct {
	Mode     string        `yaml:"mode"` // acme | selfsigned | disabled
	Email    string        `yaml:"email"`
	Resolver string        `yaml:"resolver"`
	DNS      *DNSChallenge `yaml:"dnsChallenge,omitempty"`
}

// TraefikConfig configura o Traefik
type TraefikConfig struct {
	Image       string                       `yaml:"image,omitempty"`
	Commands    []string                     `yaml:"commands,omitempty"`
	Labels      map[string]string            `yaml:"labels,omitempty"`
	Ports       []string                     `yaml:"ports,omitempty"`
	Volumes     []string                     `yaml:"volumes,omitempty"`
	Environment map[string]string            `yaml:"environment,omitempty"`
	Middlewares map[string]TraefikMiddleware `yaml:"middlewares,omitempty"`
	Plugins     map[string]TraefikPlugin     `yaml:"plugins,omitempty"`
	EntryPoints map[string]TraefikEntryPoint `yaml:"entrypoints,omitempty"`
	Providers   map[string]TraefikProvider   `yaml:"providers,omitempty"`
	API         *TraefikAPI                  `yaml:"api,omitempty"`
	Log         *TraefikLog                  `yaml:"log,omitempty"`
	AccessLog   *TraefikAccessLog            `yaml:"accessLog,omitempty"`
	Metrics     *TraefikMetrics              `yaml:"metrics,omitempty"`
}

// TraefikMiddleware define um middleware customizado
type TraefikMiddleware struct {
	AddPrefix        *MiddlewareAddPrefix        `yaml:"addPrefix,omitempty"`
	StripPrefix      *MiddlewareStripPrefix      `yaml:"stripPrefix,omitempty"`
	ReplacePathRegex *MiddlewareReplacePathRegex `yaml:"replacePathRegex,omitempty"`
	Auth             *MiddlewareAuth             `yaml:"auth,omitempty"`
	Headers          *MiddlewareHeaders          `yaml:"headers,omitempty"`
	RateLimit        *MiddlewareRateLimit        `yaml:"rateLimit,omitempty"`
	Retry            *MiddlewareRetry            `yaml:"retry,omitempty"`
	CircuitBreaker   *MiddlewareCircuitBreaker   `yaml:"circuitBreaker,omitempty"`
	Compress         *MiddlewareCompress         `yaml:"compress,omitempty"`
	CORS             *MiddlewareCORS             `yaml:"cors,omitempty"`
	CustomMiddleware map[string]interface{}      `yaml:",inline"`
}

// Middleware components
type MiddlewareAddPrefix struct {
	Prefix string `yaml:"prefix"`
}

type MiddlewareStripPrefix struct {
	Prefixes   []string `yaml:"prefixes,omitempty"`
	ForceSlash bool     `yaml:"forceSlash,omitempty"`
}

type MiddlewareReplacePathRegex struct {
	Regex       string `yaml:"regex"`
	Replacement string `yaml:"replacement"`
}

type MiddlewareAuth struct {
	Basic   *AuthBasic   `yaml:"basic,omitempty"`
	Digest  *AuthDigest  `yaml:"digest,omitempty"`
	Forward *AuthForward `yaml:"forward,omitempty"`
}

type AuthBasic struct {
	Users        []string `yaml:"users,omitempty"`
	UsersFile    string   `yaml:"usersFile,omitempty"`
	Realm        string   `yaml:"realm,omitempty"`
	RemoveHeader bool     `yaml:"removeHeader,omitempty"`
}

type AuthDigest struct {
	Users        []string `yaml:"users,omitempty"`
	UsersFile    string   `yaml:"usersFile,omitempty"`
	Realm        string   `yaml:"realm,omitempty"`
	RemoveHeader bool     `yaml:"removeHeader,omitempty"`
}

type AuthForward struct {
	Address             string          `yaml:"address"`
	TLS                 *ForwardAuthTLS `yaml:"tls,omitempty"`
	TrustForwardHeader  bool            `yaml:"trustForwardHeader,omitempty"`
	AuthResponseHeaders []string        `yaml:"authResponseHeaders,omitempty"`
}

type ForwardAuthTLS struct {
	CA                 string `yaml:"ca,omitempty"`
	CAOptional         bool   `yaml:"caOptional,omitempty"`
	Cert               string `yaml:"cert,omitempty"`
	Key                string `yaml:"key,omitempty"`
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify,omitempty"`
}

type MiddlewareHeaders struct {
	CustomRequestHeaders          map[string]string `yaml:"customRequestHeaders,omitempty"`
	CustomResponseHeaders         map[string]string `yaml:"customResponseHeaders,omitempty"`
	AccessControlAllowCredentials bool              `yaml:"accessControlAllowCredentials,omitempty"`
	AccessControlAllowHeaders     []string          `yaml:"accessControlAllowHeaders,omitempty"`
	AccessControlAllowMethods     []string          `yaml:"accessControlAllowMethods,omitempty"`
	AccessControlAllowOriginList  []string          `yaml:"accessControlAllowOriginList,omitempty"`
	AccessControlExposeHeaders    []string          `yaml:"accessControlExposeHeaders,omitempty"`
	AccessControlMaxAge           int64             `yaml:"accessControlMaxAge,omitempty"`
	AddVaryHeader                 bool              `yaml:"addVaryHeader,omitempty"`
	AllowedHosts                  []string          `yaml:"allowedHosts,omitempty"`
	BrowserXssFilter              bool              `yaml:"browserXssFilter,omitempty"`
	ContentSecurityPolicy         string            `yaml:"contentSecurityPolicy,omitempty"`
	ContentTypeNosniff            bool              `yaml:"contentTypeNosniff,omitempty"`
	ForceSTSHeader                bool              `yaml:"forceSTSHeader,omitempty"`
	FrameDeny                     bool              `yaml:"frameDeny,omitempty"`
	HostsProxyHeaders             []string          `yaml:"hostsProxyHeaders,omitempty"`
	IsDevelopment                 bool              `yaml:"isDevelopment,omitempty"`
	PublicKey                     string            `yaml:"publicKey,omitempty"`
	ReferrerPolicy                string            `yaml:"referrerPolicy,omitempty"`
	SSLForceHost                  bool              `yaml:"sslForceHost,omitempty"`
	SSLHost                       string            `yaml:"sslHost,omitempty"`
	SSLProxyHeaders               map[string]string `yaml:"sslProxyHeaders,omitempty"`
	SSLRedirect                   bool              `yaml:"sslRedirect,omitempty"`
	SSLTemporaryRedirect          bool              `yaml:"sslTemporaryRedirect,omitempty"`
	STSIncludeSubdomains          bool              `yaml:"stsIncludeSubdomains,omitempty"`
	STSPreload                    bool              `yaml:"stsPreload,omitempty"`
	STSSeconds                    int64             `yaml:"stsSeconds,omitempty"`
}

type MiddlewareRateLimit struct {
	Average         int64                     `yaml:"average,omitempty"`
	Period          string                    `yaml:"period,omitempty"`
	Burst           int64                     `yaml:"burst,omitempty"`
	SourceCriterion *RateLimitSourceCriterion `yaml:"sourceCriterion,omitempty"`
}

type RateLimitSourceCriterion struct {
	IPStrategy        *IPStrategy `yaml:"ipStrategy,omitempty"`
	RequestHeaderName string      `yaml:"requestHeaderName,omitempty"`
	RequestHost       bool        `yaml:"requestHost,omitempty"`
}

type IPStrategy struct {
	Depth       int      `yaml:"depth,omitempty"`
	ExcludedIPs []string `yaml:"excludedIPs,omitempty"`
}

type MiddlewareRetry struct {
	Attempts        int    `yaml:"attempts,omitempty"`
	InitialInterval string `yaml:"initialInterval,omitempty"`
}

type MiddlewareCircuitBreaker struct {
	Expression string `yaml:"expression"`
}

type MiddlewareCompress struct {
	ExcludedContentTypes []string `yaml:"excludedContentTypes,omitempty"`
	MinResponseBodyBytes int      `yaml:"minResponseBodyBytes,omitempty"`
}

type MiddlewareCORS struct {
	AccessControlAllowCredentials bool     `yaml:"accessControlAllowCredentials,omitempty"`
	AccessControlAllowHeaders     []string `yaml:"accessControlAllowHeaders,omitempty"`
	AccessControlAllowMethods     []string `yaml:"accessControlAllowMethods,omitempty"`
	AccessControlAllowOriginList  []string `yaml:"accessControlAllowOriginList,omitempty"`
	AccessControlExposeHeaders    []string `yaml:"accessControlExposeHeaders,omitempty"`
	AccessControlMaxAge           int64    `yaml:"accessControlMaxAge,omitempty"`
	AddVaryHeader                 bool     `yaml:"addVaryHeader,omitempty"`
}

// TraefikPlugin define um plugin customizado
type TraefikPlugin struct {
	ModuleName string                 `yaml:"moduleName"`
	Version    string                 `yaml:"version,omitempty"`
	Settings   map[string]interface{} `yaml:"settings,omitempty"`
}

// TraefikEntryPoint define um entry point customizado
type TraefikEntryPoint struct {
	Address       string                   `yaml:"address"`
	AsDefault     bool                     `yaml:"asDefault,omitempty"`
	HTTP          *EntryPointHTTP          `yaml:"http,omitempty"`
	Transport     *EntryPointTransport     `yaml:"transport,omitempty"`
	ProxyProtocol *EntryPointProxyProtocol `yaml:"proxyProtocol,omitempty"`
}

type EntryPointHTTP struct {
	Redirections *HTTPRedirections `yaml:"redirections,omitempty"`
	Middlewares  []string          `yaml:"middlewares,omitempty"`
	TLS          *EntryPointTLS    `yaml:"tls,omitempty"`
}

type HTTPRedirections struct {
	EntryPoint *HTTPRedirectionEntryPoint `yaml:"entryPoint,omitempty"`
}

type HTTPRedirectionEntryPoint struct {
	To        string `yaml:"to,omitempty"`
	Scheme    string `yaml:"scheme,omitempty"`
	Permanent bool   `yaml:"permanent,omitempty"`
}

type EntryPointTLS struct {
	Options      string   `yaml:"options,omitempty"`
	CertResolver string   `yaml:"certResolver,omitempty"`
	Domains      []Domain `yaml:"domains,omitempty"`
}

type Domain struct {
	Main string   `yaml:"main"`
	SANs []string `yaml:"sans,omitempty"`
}

type EntryPointTransport struct {
	RespondingTimeouts   *RespondingTimeouts `yaml:"respondingTimeouts,omitempty"`
	LifeCycle            *LifeCycle          `yaml:"lifeCycle,omitempty"`
	KeepAliveMaxRequests int                 `yaml:"keepAliveMaxRequests,omitempty"`
	KeepAliveMaxTime     string              `yaml:"keepAliveMaxTime,omitempty"`
}

type RespondingTimeouts struct {
	ReadTimeout  string `yaml:"readTimeout,omitempty"`
	WriteTimeout string `yaml:"writeTimeout,omitempty"`
	IdleTimeout  string `yaml:"idleTimeout,omitempty"`
}

type LifeCycle struct {
	RequestAcceptGraceTimeout string `yaml:"requestAcceptGraceTimeout,omitempty"`
	GraceTimeOut              string `yaml:"graceTimeOut,omitempty"`
}

type EntryPointProxyProtocol struct {
	Insecure   bool     `yaml:"insecure,omitempty"`
	TrustedIPs []string `yaml:"trustedIPs,omitempty"`
}

// TraefikProvider define configurações de providers
type TraefikProvider struct {
	Docker     *DockerProvider     `yaml:"docker,omitempty"`
	File       *FileProvider       `yaml:"file,omitempty"`
	Consul     *ConsulProvider     `yaml:"consul,omitempty"`
	Kubernetes *KubernetesProvider `yaml:"kubernetes,omitempty"`
}

type DockerProvider struct {
	Constraints             string `yaml:"constraints,omitempty"`
	Watch                   bool   `yaml:"watch,omitempty"`
	Endpoint                string `yaml:"endpoint,omitempty"`
	DefaultRule             string `yaml:"defaultRule,omitempty"`
	ExposedByDefault        bool   `yaml:"exposedByDefault,omitempty"`
	UseBindPortIP           bool   `yaml:"useBindPortIP,omitempty"`
	SwarmMode               bool   `yaml:"swarmMode,omitempty"`
	Network                 string `yaml:"network,omitempty"`
	SwarmModeRefreshSeconds int    `yaml:"swarmModeRefreshSeconds,omitempty"`
}

type FileProvider struct {
	Directory                 string `yaml:"directory,omitempty"`
	Watch                     bool   `yaml:"watch,omitempty"`
	Filename                  string `yaml:"filename,omitempty"`
	DebugLogGeneratedTemplate bool   `yaml:"debugLogGeneratedTemplate,omitempty"`
}

type ConsulProvider struct {
	Endpoints []string `yaml:"endpoints,omitempty"`
	RootKey   string   `yaml:"rootKey,omitempty"`
}

type KubernetesProvider struct {
	Endpoint         string           `yaml:"endpoint,omitempty"`
	Token            string           `yaml:"token,omitempty"`
	CertAuthFilePath string           `yaml:"certAuthFilePath,omitempty"`
	Namespaces       []string         `yaml:"namespaces,omitempty"`
	LabelSelector    string           `yaml:"labelSelector,omitempty"`
	IngressClass     string           `yaml:"ingressClass,omitempty"`
	IngressEndpoint  *IngressEndpoint `yaml:"ingressEndpoint,omitempty"`
	ThrottleDuration string           `yaml:"throttleDuration,omitempty"`
}

type IngressEndpoint struct {
	IP               string `yaml:"ip,omitempty"`
	Hostname         string `yaml:"hostname,omitempty"`
	PublishedService string `yaml:"publishedService,omitempty"`
}

// TraefikAPI define configurações da API
type TraefikAPI struct {
	Dashboard     bool   `yaml:"dashboard,omitempty"`
	Debug         bool   `yaml:"debug,omitempty"`
	Insecure      bool   `yaml:"insecure,omitempty"`
	DashboardPath string `yaml:"dashboardPath,omitempty"`
}

// TraefikLog define configurações de log
type TraefikLog struct {
	Level    string `yaml:"level,omitempty"`
	Format   string `yaml:"format,omitempty"`
	FilePath string `yaml:"filePath,omitempty"`
}

// TraefikAccessLog define configurações de access log
type TraefikAccessLog struct {
	FilePath string            `yaml:"filePath,omitempty"`
	Format   string            `yaml:"format,omitempty"`
	Filters  *AccessLogFilters `yaml:"filters,omitempty"`
	Fields   *AccessLogFields  `yaml:"fields,omitempty"`
}

type AccessLogFilters struct {
	StatusCodes   []string `yaml:"statusCodes,omitempty"`
	RetryAttempts bool     `yaml:"retryAttempts,omitempty"`
	MinDuration   string   `yaml:"minDuration,omitempty"`
}

type AccessLogFields struct {
	DefaultMode string                  `yaml:"defaultMode,omitempty"`
	Names       map[string]string       `yaml:"names,omitempty"`
	Headers     *AccessLogFieldsHeaders `yaml:"headers,omitempty"`
}

type AccessLogFieldsHeaders struct {
	DefaultMode string            `yaml:"defaultMode,omitempty"`
	Names       map[string]string `yaml:"names,omitempty"`
}

// TraefikMetrics define configurações de métricas
type TraefikMetrics struct {
	Prometheus *PrometheusMetrics `yaml:"prometheus,omitempty"`
	Datadog    *DatadogMetrics    `yaml:"datadog,omitempty"`
	StatsD     *StatsDMetrics     `yaml:"statsD,omitempty"`
	InfluxDB   *InfluxDBMetrics   `yaml:"influxDB,omitempty"`
}

type PrometheusMetrics struct {
	AddEntryPointsLabels bool      `yaml:"addEntryPointsLabels,omitempty"`
	AddServicesLabels    bool      `yaml:"addServicesLabels,omitempty"`
	Buckets              []float64 `yaml:"buckets,omitempty"`
}

type DatadogMetrics struct {
	Address      string `yaml:"address,omitempty"`
	PushInterval string `yaml:"pushInterval,omitempty"`
}

type StatsDMetrics struct {
	Address      string `yaml:"address,omitempty"`
	PushInterval string `yaml:"pushInterval,omitempty"`
}

type InfluxDBMetrics struct {
	Address         string `yaml:"address,omitempty"`
	Protocol        string `yaml:"protocol,omitempty"`
	PushInterval    string `yaml:"pushInterval,omitempty"`
	Database        string `yaml:"database,omitempty"`
	RetentionPolicy string `yaml:"retentionPolicy,omitempty"`
	Username        string `yaml:"username,omitempty"`
	Password        string `yaml:"password,omitempty"`
}

// DNSChallenge configura DNS challenge
type DNSChallenge struct {
	Provider string   `yaml:"provider"`
	Env      []string `yaml:"env"`
}

// Observability configura monitoramento
type Observability struct {
	Dozzle       Dozzle `yaml:"dozzle"`
	Beszel       Beszel `yaml:"beszel"`
	DockerSocket string `yaml:"docker_socket,omitempty"` // Path customizado do Docker socket
}

// Dozzle configura logs
type Dozzle struct {
	Enabled    bool       `yaml:"enabled"`
	Subdomain  string     `yaml:"subdomain"`
	DataVolume string     `yaml:"data_volume"`
	BasicAuth  *BasicAuth `yaml:"basic_auth,omitempty"`
}

// Beszel configura monitoramento via socket Unix
type Beszel struct {
	Enabled      bool   `yaml:"enabled"`
	Subdomain    string `yaml:"subdomain"`
	DataVolume   string `yaml:"data_volume"`
	SocketVolume string `yaml:"socket_volume"`

	// Configuração de autenticação - apenas o essencial
	PublicKey string `yaml:"public_key,omitempty"` // Chave pública SSH para autenticação
	Token     string `yaml:"token,omitempty"`      // Token de autenticação do agent

	// Configurações avançadas opcionais
	AppURL       string `yaml:"app_url,omitempty"`       // URL customizada do hub
	UserCreation bool   `yaml:"user_creation,omitempty"` // Permitir criação automática de usuários
}

// Network representa uma network Docker
type Network struct {
	Internal bool `yaml:"internal"`
}

// Volume representa um volume Docker
type Volume struct {
	Name string `yaml:"name"`
}

// BuildSpec configura build
type BuildSpec struct {
	Context    string            `yaml:"context"`
	Dockerfile string            `yaml:"dockerfile"`
	Args       map[string]string `yaml:"args,omitempty"`
}

// VolumeMount representa um mount de volume
type VolumeMount struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

// Service representa um serviço
type Service struct {
	Name          string            `yaml:"name"`
	Subdomain     string            `yaml:"subdomain,omitempty"`
	Image         string            `yaml:"image,omitempty"`
	Build         *BuildSpec        `yaml:"build,omitempty"`
	Expose        int               `yaml:"expose"`
	Replicas      int               `yaml:"replicas,omitempty"`
	Env           map[string]string `yaml:"env,omitempty"`
	EnvFile       []string          `yaml:"env_file,omitempty"`
	Secrets       []Secret          `yaml:"secrets,omitempty"`
	Volumes       []VolumeMount     `yaml:"volumes,omitempty"`
	Resources     *Resources        `yaml:"resources,omitempty"`
	HealthCheck   *HealthCheck      `yaml:"health_check,omitempty"`
	Deploy        *DeployConfig     `yaml:"deploy,omitempty"`
	TraefikRaw    interface{}       `yaml:"traefik,omitempty"`
	BasicAuth     *BasicAuth        `yaml:"basic_auth,omitempty"`
	NetworkAccess *NetworkAccess    `yaml:"network_access,omitempty"`
}

// GetTraefik retorna a configuração do Traefik com compatibilidade para ambos os formatos
func (s *Service) GetTraefik() *ServiceTraefik {
	if s.TraefikRaw == nil {
		return nil
	}

	switch v := s.TraefikRaw.(type) {
	case bool:
		// Formato antigo: traefik: true/false
		if v {
			return &ServiceTraefik{Enabled: true}
		}
		return &ServiceTraefik{Enabled: false}
	case map[string]interface{}:
		// Formato novo: estrutura complexa
		traefik := &ServiceTraefik{}

		if enabled, ok := v["enabled"].(bool); ok {
			traefik.Enabled = enabled
		}

		if rule, ok := v["rule"].(string); ok {
			traefik.Rule = rule
		}

		if priority, ok := v["priority"].(int); ok {
			traefik.Priority = priority
		}

		if entrypoints, ok := v["entrypoints"].([]interface{}); ok {
			for _, ep := range entrypoints {
				if epStr, ok := ep.(string); ok {
					traefik.EntryPoints = append(traefik.EntryPoints, epStr)
				}
			}
		}

		if middlewares, ok := v["middlewares"].([]interface{}); ok {
			for _, mw := range middlewares {
				if mwStr, ok := mw.(string); ok {
					traefik.Middlewares = append(traefik.Middlewares, mwStr)
				}
			}
		}

		if labels, ok := v["labels"].(map[string]interface{}); ok {
			traefik.Labels = make(map[string]string)
			for k, label := range labels {
				if labelStr, ok := label.(string); ok {
					traefik.Labels[k] = labelStr
				}
			}
		}

		return traefik
	case *ServiceTraefik:
		// Formato já parseado
		return v
	default:
		return nil
	}
}

// ServiceTraefik define configurações específicas do Traefik para um serviço
type ServiceTraefik struct {
	Enabled      bool              `yaml:"enabled,omitempty"`
	Rule         string            `yaml:"rule,omitempty"`         // Regra customizada (sobrescreve padrão Host())
	EntryPoints  []string          `yaml:"entrypoints,omitempty"`  // Entry points customizados
	Middlewares  []string          `yaml:"middlewares,omitempty"`  // Lista de middlewares a aplicar
	Priority     int               `yaml:"priority,omitempty"`     // Prioridade da rota
	TLS          *ServiceTLS       `yaml:"tls,omitempty"`          // Configurações TLS específicas
	LoadBalancer *ServiceLB        `yaml:"loadBalancer,omitempty"` // Configurações load balancer
	Service      *ServiceConfig    `yaml:"service,omitempty"`      // Configurações avançadas do serviço
	Labels       map[string]string `yaml:"labels,omitempty"`       // Labels Traefik customizados
}

// ServiceTLS define configurações TLS específicas do serviço
type ServiceTLS struct {
	CertResolver string   `yaml:"certResolver,omitempty"`
	Domains      []Domain `yaml:"domains,omitempty"`
	Options      string   `yaml:"options,omitempty"`
	Store        string   `yaml:"store,omitempty"`
}

// ServiceLB define configurações de load balancer
type ServiceLB struct {
	Sticky             *StickyConfig       `yaml:"sticky,omitempty"`
	HealthCheck        *LBHealthCheck      `yaml:"healthCheck,omitempty"`
	PassHostHeader     bool                `yaml:"passHostHeader,omitempty"`
	ResponseForwarding *ResponseForwarding `yaml:"responseForwarding,omitempty"`
	ServersTransport   string              `yaml:"serversTransport,omitempty"`
}

// StickyConfig define configurações de sessão sticky
type StickyConfig struct {
	Cookie *StickyCookie `yaml:"cookie,omitempty"`
}

type StickyCookie struct {
	Name     string `yaml:"name,omitempty"`
	Secure   bool   `yaml:"secure,omitempty"`
	HTTPOnly bool   `yaml:"httpOnly,omitempty"`
	SameSite string `yaml:"sameSite,omitempty"`
}

// LBHealthCheck define configurações de health check do load balancer
type LBHealthCheck struct {
	Path            string            `yaml:"path,omitempty"`
	Port            int               `yaml:"port,omitempty"`
	Interval        string            `yaml:"interval,omitempty"`
	Timeout         string            `yaml:"timeout,omitempty"`
	Hostname        string            `yaml:"hostname,omitempty"`
	FollowRedirects bool              `yaml:"followRedirects,omitempty"`
	Headers         map[string]string `yaml:"headers,omitempty"`
	Method          string            `yaml:"method,omitempty"`
	Status          int               `yaml:"status,omitempty"`
	Scheme          string            `yaml:"scheme,omitempty"`
}

// ResponseForwarding define configurações de resposta
type ResponseForwarding struct {
	FlushInterval string `yaml:"flushInterval,omitempty"`
}

// ServiceConfig define configurações avançadas do serviço Traefik
type ServiceConfig struct {
	PassTLSCert bool `yaml:"passTLSCert,omitempty"`
}

// Secret representa uma secret Docker
type Secret struct {
	Name     string `yaml:"name"`
	File     string `yaml:"file,omitempty"`
	External bool   `yaml:"external,omitempty"`
	Target   string `yaml:"target,omitempty"`
}

// Resources define limites de recursos
type Resources struct {
	Memory     string            `yaml:"memory,omitempty"`   // ex: "512m", "1g"
	CPUs       string            `yaml:"cpus,omitempty"`     // ex: "0.5", "1.0"
	GPUs       string            `yaml:"gpus,omitempty"`     // ex: "1", "all"
	ShmSize    string            `yaml:"shm_size,omitempty"` // ex: "128m"
	Ulimits    map[string]Ulimit `yaml:"ulimits,omitempty"`
	ReserveCPU string            `yaml:"reserve_cpu,omitempty"` // reserva mínima
	ReserveMem string            `yaml:"reserve_mem,omitempty"` // reserva mínima
}

// Ulimit representa um ulimit
type Ulimit struct {
	Soft int64 `yaml:"soft"`
	Hard int64 `yaml:"hard"`
}

// BasicAuth representa as configurações de autenticação básica
type BasicAuth struct {
	Enabled   bool              `yaml:"enabled"`
	Username  string            `yaml:"username,omitempty"`
	Password  string            `yaml:"password,omitempty"`
	Users     map[string]string `yaml:"users,omitempty"`      // user:password
	UsersFile string            `yaml:"users_file,omitempty"` // arquivo htpasswd
}

// HealthCheck representa as configurações de health check
type HealthCheck struct {
	Enabled  bool   `yaml:"enabled"`
	Path     string `yaml:"path,omitempty"`
	Interval string `yaml:"interval,omitempty"`
	Timeout  string `yaml:"timeout,omitempty"`
	Retries  int    `yaml:"retries,omitempty"`
}

// DeployConfig representa as configurações de deployment
type DeployConfig struct {
	Strategy string `yaml:"strategy,omitempty"` // "rolling" ou "recreate"
}

// NetworkAccess define o acesso às redes
type NetworkAccess struct {
	Internet bool     `yaml:"internet,omitempty"` // Permite acesso à internet (rede pública)
	Internal bool     `yaml:"internal,omitempty"` // Força apenas rede privada (padrão: true)
	Custom   []string `yaml:"custom,omitempty"`   // Redes customizadas adicionais
}

// applyDefaults aplica valores padrão
func (s *Stack) applyDefaults() {
	// defaults de observability
	if s.Observability.Dozzle.DataVolume == "" {
		s.Observability.Dozzle.DataVolume = "dozzle_data"
	}
	if s.Observability.Beszel.DataVolume == "" {
		s.Observability.Beszel.DataVolume = "beszel_data"
	}
	if s.Observability.Beszel.SocketVolume == "" {
		s.Observability.Beszel.SocketVolume = "beszel_socket"
	}

	// habilita por padrão se nenhum estiver habilitado
	if !s.Observability.Dozzle.Enabled && !s.Observability.Beszel.Enabled {
		s.Observability.Dozzle.Enabled = true
		s.Observability.Beszel.Enabled = true
	}

	if s.TLS.Resolver == "" {
		s.TLS.Resolver = "le"
	}
}
