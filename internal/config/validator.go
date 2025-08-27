package config

import (
	"context"
	"errors"
	"fmt"
)

// Validator validates configurations
type Validator interface {
	Validate(ctx context.Context, stack *Stack) error
}

// validator implements Validator
type validator struct{}

// NewValidator creates a new validator
func NewValidator() Validator {
	return &validator{}
}

func (v *validator) Validate(ctx context.Context, stack *Stack) error {
	stack.applyDefaults()

	var errs []error

	// Basic validations
	if stack.Version != 1 {
		errs = append(errs, errors.New("version must be 1"))
	}
	if stack.Project == "" {
		errs = append(errs, errors.New("project is required"))
	}
	if stack.Domain == "" {
		errs = append(errs, errors.New("domain is required"))
	}

	// TLS validation
	if err := v.validateTLS(&stack.TLS); err != nil {
		errs = append(errs, err)
	}

	// Networks validation
	if err := v.validateNetworks(stack.Networks); err != nil {
		errs = append(errs, err)
	}

	// Services validation
	if err := v.validateServices(stack.Services); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		msg := "invalid config:\n"
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
		// valid modes
	default:
		return fmt.Errorf("invalid tls.mode: %q", tls.Mode)
	}

	if tls.Mode == "acme" && tls.Email == "" {
		return errors.New("tls.email required with acme")
	}

	return nil
}

func (v *validator) validateNetworks(networks map[string]Network) error {
	if _, ok := networks["public"]; !ok {
		return errors.New("network 'public' is required")
	}
	if _, ok := networks["private"]; !ok {
		return errors.New("network 'private' is required")
	}
	return nil
}

func (v *validator) validateServices(services []Service) error {
	if len(services) == 0 {
		return errors.New("define at least one service")
	}

	seen := make(map[string]struct{})
	for _, sv := range services {
		if sv.Name == "" {
			return errors.New("service.name is required")
		}

		if _, ok := seen[sv.Name]; ok {
			return fmt.Errorf("duplicate service: %s", sv.Name)
		}
		seen[sv.Name] = struct{}{}

		if sv.Image == "" && sv.Build == nil {
			return fmt.Errorf("%s: define image OR build", sv.Name)
		}
		if sv.Image != "" && sv.Build != nil {
			return fmt.Errorf("%s: use image OR build (not both)", sv.Name)
		}
		if sv.Expose <= 0 {
			return fmt.Errorf("%s: expose must be > 0", sv.Name)
		}
		if sv.Traefik && sv.Subdomain == "" {
			return fmt.Errorf("%s: subdomain is required when traefik=true", sv.Name)
		}

		// Validate replicas
		if sv.Replicas < 0 {
			return fmt.Errorf("%s: replicas cannot be negative", sv.Name)
		}

		// Validate volumes
		for _, m := range sv.Volumes {
			if m.Source == "" || m.Target == "" {
				return fmt.Errorf("%s: invalid volume (source/target)", sv.Name)
			}
		}

		// Validate secrets
		for _, secret := range sv.Secrets {
			if secret.Name == "" {
				return fmt.Errorf("%s: secret.name is required", sv.Name)
			}
			if !secret.External && secret.File == "" {
				return fmt.Errorf("%s: secret '%s' needs 'file' or 'external=true'", sv.Name, secret.Name)
			}
		}

		// Validate basic auth
		if sv.BasicAuth != nil && sv.BasicAuth.Enabled {
			if len(sv.BasicAuth.Users) == 0 && sv.BasicAuth.UsersFile == "" {
				return fmt.Errorf("%s: basic_auth enabled needs 'users' or 'users_file'", sv.Name)
			}
		}

		// Validate resources
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

	// Validate memory format
	if resources.Memory != "" {
		if err := v.validateMemoryFormat(resources.Memory); err != nil {
			return fmt.Errorf("%s: invalid memory: %v", serviceName, err)
		}
	}

	// Validate CPU format
	if resources.CPUs != "" {
		if err := v.validateCPUFormat(resources.CPUs); err != nil {
			return fmt.Errorf("%s: invalid cpus: %v", serviceName, err)
		}
	}

	return nil
}

func (v *validator) validateMemoryFormat(memory string) error {
	// Valid formats: 512m, 1g, 2048M, 1G, etc.
	if len(memory) < 2 {
		return errors.New("invalid format")
	}

	unit := memory[len(memory)-1:]
	if unit != "m" && unit != "M" && unit != "g" && unit != "G" {
		return errors.New("unit must be m, M, g or G")
	}

	return nil
}

func (v *validator) validateCPUFormat(cpu string) error {
	// Valid formats: 0.5, 1, 1.0, 2, etc.
	if cpu == "" {
		return errors.New("empty value")
	}

	return nil
}
