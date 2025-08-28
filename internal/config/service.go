package config

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/leandrodaf/harborctl/pkg/fs"
)

// Manager gerencia configurações
type Manager interface {
	Load(ctx context.Context, path string) (*Stack, error)
	Validate(ctx context.Context, stack *Stack) error
	Create(ctx context.Context, path string, options CreateOptions) error
	SaveBaseConfig(ctx context.Context, path string, stack *Stack) error
}

// CreateOptions configura a criação de stack
type CreateOptions struct {
	Domain         string
	Email          string
	Project        string
	Environment    string
	NoDozzle       bool
	NoBeszel       bool
	DozzleAuth     bool
	BeszelAuth     bool
	DozzleUsername string
	DozzlePassword string
	BeszelUsername string
	BeszelPassword string
}

// manager implementa Manager
type manager struct {
	loader    fs.ConfigLoader
	fs        fs.FileSystem
	validator Validator
}

// NewManager cria um novo gerenciador de configuração
func NewManager(loader fs.ConfigLoader, filesystem fs.FileSystem, validator Validator) Manager {
	return &manager{
		loader:    loader,
		fs:        filesystem,
		validator: validator,
	}
}

func (m *manager) Load(ctx context.Context, path string) (*Stack, error) {
	data, err := m.loader.Load(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	var stack Stack
	if err := yaml.Unmarshal(data, &stack); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &stack, nil
}

func (m *manager) Validate(ctx context.Context, stack *Stack) error {
	return m.validator.Validate(ctx, stack)
}

func (m *manager) Create(ctx context.Context, path string, options CreateOptions) error {
	if m.fs.Exists(path) {
		return errors.New("stack.yml already exists")
	}

	// Detectar ambiente automaticamente se não especificado
	env := options.Environment
	if env == "" {
		if options.Domain == "localhost" || options.Domain == "test.local" ||
			strings.HasSuffix(options.Domain, ".local") || strings.HasSuffix(options.Domain, ".localhost") {
			env = "local"
		} else {
			env = "production"
		}
	}

	// Configuração de TLS baseada no ambiente
	var tlsConfig TLS
	if env == "local" {
		tlsConfig = TLS{
			Mode: "disabled",
		}
	} else {
		tlsConfig = TLS{
			Mode:     "acme",
			Email:    options.Email,
			Resolver: "le",
		}
	}

	var dozzleBasicAuth *BasicAuth
	if options.DozzleAuth && !options.NoDozzle {
		dozzleBasicAuth = &BasicAuth{
			Enabled:  true,
			Username: options.DozzleUsername,
			Password: options.DozzlePassword,
		}
	}

	var beszelBasicAuth *BasicAuth
	if options.BeszelAuth && !options.NoBeszel {
		beszelBasicAuth = &BasicAuth{
			Enabled:  true,
			Username: options.BeszelUsername,
			Password: options.BeszelPassword,
		}
	}

	stack := &Stack{
		Version:     1,
		Project:     options.Project,
		Domain:      options.Domain,
		Environment: env,
		TLS:         tlsConfig,
		Observability: Observability{
			Dozzle: Dozzle{
				Enabled:    !options.NoDozzle,
				Subdomain:  "logs",
				DataVolume: "dozzle_data",
				BasicAuth:  dozzleBasicAuth,
			},
			Beszel: Beszel{
				Enabled:      !options.NoBeszel,
				Subdomain:    "monitor",
				DataVolume:   "beszel_data",
				SocketVolume: "beszel_socket",
				BasicAuth:    beszelBasicAuth,
			},
		},
		Networks: map[string]Network{
			"private": {Internal: true},
			"public":  {Internal: false},
		},
		Volumes: []Volume{
			{Name: "traefik_acme"},
			{Name: "dozzle_data"},
			{Name: "beszel_data"},
			{Name: "beszel_socket"},
		},
		Services: []Service{},
	}

	data, err := yaml.Marshal(stack)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return m.fs.WriteFile(path, data, 0644)
}

// SaveBaseConfig salva a configuração base do servidor
func (m *manager) SaveBaseConfig(ctx context.Context, path string, stack *Stack) error {
	data, err := yaml.Marshal(stack)
	if err != nil {
		return fmt.Errorf("failed to marshal base config: %w", err)
	}

	return m.fs.WriteFile(path, data, 0644)
}
