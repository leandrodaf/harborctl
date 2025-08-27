package config

// Stack representa a configuração completa
type Stack struct {
	Version       int                `yaml:"version"`
	Project       string             `yaml:"project"`
	Domain        string             `yaml:"domain"`
	TLS           TLS                `yaml:"tls"`
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

// DNSChallenge configura DNS challenge
type DNSChallenge struct {
	Provider string   `yaml:"provider"`
	Env      []string `yaml:"env"`
}

// Observability configura monitoramento
type Observability struct {
	Dozzle Dozzle `yaml:"dozzle"`
	Beszel Beszel `yaml:"beszel"`
}

// Dozzle configura logs
type Dozzle struct {
	Enabled    bool   `yaml:"enabled"`
	Subdomain  string `yaml:"subdomain"`
	DataVolume string `yaml:"data_volume"`
}

// Beszel configura monitoramento
type Beszel struct {
	Enabled      bool   `yaml:"enabled"`
	Subdomain    string `yaml:"subdomain"`
	DataVolume   string `yaml:"data_volume"`
	SocketVolume string `yaml:"socket_volume"`
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
	Name        string            `yaml:"name"`
	Subdomain   string            `yaml:"subdomain,omitempty"`
	Image       string            `yaml:"image,omitempty"`
	Build       *BuildSpec        `yaml:"build,omitempty"`
	Expose      int               `yaml:"expose"`
	Replicas    int               `yaml:"replicas,omitempty"`
	Env         map[string]string `yaml:"env,omitempty"`
	EnvFile     []string          `yaml:"env_file,omitempty"`
	Secrets     []Secret          `yaml:"secrets,omitempty"`
	Volumes     []VolumeMount     `yaml:"volumes,omitempty"`
	Resources   *Resources        `yaml:"resources,omitempty"`
	HealthCheck *HealthCheck      `yaml:"health_check,omitempty"`
	Deploy      *DeployConfig     `yaml:"deploy,omitempty"`
	Traefik     bool              `yaml:"traefik"`
	BasicAuth   *BasicAuth        `yaml:"basic_auth,omitempty"`
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
