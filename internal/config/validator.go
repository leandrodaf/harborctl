package config

import (
	"context"
	"errors"
	"fmt"
)

// Validator valida configurações
type Validator interface {
	Validate(ctx context.Context, stack *Stack) error
}

// validator implementa Validator
type validator struct{}

// NewValidator cria um novo validator
func NewValidator() Validator {
	return &validator{}
}

func (v *validator) Validate(ctx context.Context, stack *Stack) error {
	stack.applyDefaults()

	var errs []error

	// Validações básicas
	if stack.Version != 1 {
		errs = append(errs, errors.New("version deve ser 1"))
	}
	if stack.Project == "" {
		errs = append(errs, errors.New("project é obrigatório"))
	}
	if stack.Domain == "" {
		errs = append(errs, errors.New("domain é obrigatório"))
	}

	// Validação do TLS
	if err := v.validateTLS(&stack.TLS); err != nil {
		errs = append(errs, err)
	}

	// Validação das networks
	if err := v.validateNetworks(stack.Networks); err != nil {
		errs = append(errs, err)
	}

	// Validação dos services
	if err := v.validateServices(stack.Services); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		msg := "config inválida:\n"
		for _, e := range errs {
			msg += " - " + e.Error() + "\n"
		}
		return errors.New(msg)
	}

	return nil
}

func (v *validator) validateTLS(tls *TLS) error {
	switch tls.Mode {
	case "acme", "selfsigned", "disabled":
		// modos válidos
	default:
		return fmt.Errorf("tls.mode inválido: %q", tls.Mode)
	}

	if tls.Mode == "acme" && tls.Email == "" {
		return errors.New("tls.email obrigatório com acme")
	}

	return nil
}

func (v *validator) validateNetworks(networks map[string]Network) error {
	if _, ok := networks["public"]; !ok {
		return errors.New("network 'public' é obrigatória")
	}
	if _, ok := networks["private"]; !ok {
		return errors.New("network 'private' é obrigatória")
	}
	return nil
}

func (v *validator) validateServices(services []Service) error {
	if len(services) == 0 {
		return errors.New("defina ao menos um service")
	}

	seen := make(map[string]struct{})
	for _, sv := range services {
		if sv.Name == "" {
			return errors.New("service.name é obrigatório")
		}

		if _, ok := seen[sv.Name]; ok {
			return fmt.Errorf("service duplicado: %s", sv.Name)
		}
		seen[sv.Name] = struct{}{}

		if sv.Image == "" && sv.Build == nil {
			return fmt.Errorf("%s: defina image OU build", sv.Name)
		}
		if sv.Image != "" && sv.Build != nil {
			return fmt.Errorf("%s: use image OU build (não ambos)", sv.Name)
		}
		if sv.Expose <= 0 {
			return fmt.Errorf("%s: expose deve ser > 0", sv.Name)
		}
		if sv.Traefik && sv.Subdomain == "" {
			return fmt.Errorf("%s: subdomain é obrigatório quando traefik=true", sv.Name)
		}

		// Validar replicas
		if sv.Replicas < 0 {
			return fmt.Errorf("%s: replicas não pode ser negativo", sv.Name)
		}

		// Validar volumes
		for _, m := range sv.Volumes {
			if m.Source == "" || m.Target == "" {
				return fmt.Errorf("%s: volume inválido (source/target)", sv.Name)
			}
		}

		// Validar secrets
		for _, secret := range sv.Secrets {
			if secret.Name == "" {
				return fmt.Errorf("%s: secret.name é obrigatório", sv.Name)
			}
			if !secret.External && secret.File == "" {
				return fmt.Errorf("%s: secret '%s' precisa de 'file' ou 'external=true'", sv.Name, secret.Name)
			}
		}

		// Validar basic auth
		if sv.BasicAuth != nil && sv.BasicAuth.Enabled {
			if len(sv.BasicAuth.Users) == 0 && sv.BasicAuth.UsersFile == "" {
				return fmt.Errorf("%s: basic_auth habilitado precisa de 'users' ou 'users_file'", sv.Name)
			}
		}

		// Validar recursos
		if err := v.validateResources(sv.Name, sv.Resources); err != nil {
			return err
		}
	}

	return nil
}

func (v *validator) validateResources(serviceName string, resources *Resources) error {
	if resources == nil {
		return nil
	}

	// Validar formato de memória
	if resources.Memory != "" {
		if err := v.validateMemoryFormat(resources.Memory); err != nil {
			return fmt.Errorf("%s: memory inválida: %v", serviceName, err)
		}
	}

	// Validar formato de CPU
	if resources.CPUs != "" {
		if err := v.validateCPUFormat(resources.CPUs); err != nil {
			return fmt.Errorf("%s: cpus inválido: %v", serviceName, err)
		}
	}

	return nil
}

func (v *validator) validateMemoryFormat(memory string) error {
	// Formatos válidos: 512m, 1g, 2048M, 1G, etc.
	if len(memory) < 2 {
		return errors.New("formato inválido")
	}

	unit := memory[len(memory)-1:]
	if unit != "m" && unit != "M" && unit != "g" && unit != "G" {
		return errors.New("unidade deve ser m, M, g ou G")
	}

	return nil
}

func (v *validator) validateCPUFormat(cpu string) error {
	// Formatos válidos: 0.5, 1, 1.0, 2, etc.
	if cpu == "" {
		return errors.New("valor vazio")
	}

	return nil
}
