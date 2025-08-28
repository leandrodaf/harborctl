package config

import (
	"fmt"
	"strings"

	"github.com/leandrodaf/harborctl/internal/security"
)

// SecureValidator implementa validação com segurança
type SecureValidator struct {
	pathValidator  *security.PathValidator
	inputSanitizer *security.InputSanitizer
}

// NewSecureValidator cria um novo validador seguro
func NewSecureValidator() *SecureValidator {
	return &SecureValidator{
		pathValidator:  security.NewPathValidator(),
		inputSanitizer: security.NewInputSanitizer(1000), // Max 1000 chars
	}
}

// ValidateStack valida uma stack com verificações de segurança
func (sv *SecureValidator) ValidateStack(stack *Stack) error {
	// Valida domínio
	if err := security.ValidateDomainName(stack.Domain); err != nil {
		return fmt.Errorf("domínio inválido: %w", err)
	}

	// Valida email apenas se TLS estiver habilitado e for ACME
	if stack.TLS.Mode == "acme" {
		if err := security.ValidateEmail(stack.TLS.Email); err != nil {
			return fmt.Errorf("email inválido: %w", err)
		}
	}

	// Valida serviços
	for i, service := range stack.Services {
		if err := sv.ValidateService(service); err != nil {
			return fmt.Errorf("serviço %d (%s) inválido: %w", i, service.Name, err)
		}
	}

	return nil
}

// ValidateService valida um serviço
func (sv *SecureValidator) ValidateService(service Service) error {
	// Valida nome do serviço
	if service.Name == "" {
		return fmt.Errorf("nome do serviço é obrigatório")
	}

	cleanName, err := sv.inputSanitizer.SanitizeString(service.Name)
	if err != nil {
		return fmt.Errorf("nome do serviço inválido: %w", err)
	}

	if cleanName != service.Name {
		return fmt.Errorf("nome do serviço contém caracteres inválidos")
	}

	// Valida subdomain
	if service.Subdomain != "" {
		if err := security.ValidateDomainName(service.Subdomain); err != nil {
			return fmt.Errorf("subdomínio inválido: %w", err)
		}
	}

	// Valida build context se especificado
	if service.Build != nil {
		if err := sv.ValidateBuildSpec(*service.Build); err != nil {
			return fmt.Errorf("build spec inválido: %w", err)
		}
	}

	// Valida volumes
	for i, volume := range service.Volumes {
		if err := sv.ValidateVolumeMount(volume); err != nil {
			return fmt.Errorf("volume %d inválido: %w", i, err)
		}
	}

	// Valida resources
	if service.Resources != nil {
		if err := sv.ValidateResources(*service.Resources); err != nil {
			return fmt.Errorf("recursos inválidos: %w", err)
		}
	}

	// Valida env files
	for i, envFile := range service.EnvFile {
		if err := sv.pathValidator.ValidatePath(envFile); err != nil {
			return fmt.Errorf("env file %d inválido: %w", i, err)
		}
	}

	// Valida secrets
	for i, secret := range service.Secrets {
		if err := sv.ValidateSecret(secret); err != nil {
			return fmt.Errorf("secret %d inválido: %w", i, err)
		}
	}

	return nil
}

// ValidateBuildSpec valida especificação de build
func (sv *SecureValidator) ValidateBuildSpec(build BuildSpec) error {
	// Valida context
	if err := sv.pathValidator.ValidatePath(build.Context); err != nil {
		return fmt.Errorf("build context inválido: %w", err)
	}

	// Valida dockerfile
	if build.Dockerfile != "" {
		if err := sv.pathValidator.ValidateFileName(build.Dockerfile); err != nil {
			return fmt.Errorf("dockerfile inválido: %w", err)
		}

		// Verifica se é um Dockerfile válido
		if !strings.HasSuffix(strings.ToLower(build.Dockerfile), "dockerfile") &&
			!strings.Contains(strings.ToLower(build.Dockerfile), "dockerfile") {
			return fmt.Errorf("arquivo dockerfile deve conter 'dockerfile' no nome")
		}
	}

	// Valida args
	for key, value := range build.Args {
		cleanKey, err := sv.inputSanitizer.SanitizeString(key)
		if err != nil {
			return fmt.Errorf("build arg key inválido '%s': %w", key, err)
		}

		cleanValue, err := sv.inputSanitizer.SanitizeString(value)
		if err != nil {
			return fmt.Errorf("build arg value inválido '%s': %w", value, err)
		}

		if cleanKey != key || cleanValue != value {
			return fmt.Errorf("build args contém caracteres inválidos")
		}
	}

	return nil
}

// ValidateVolumeMount valida um mount de volume
func (sv *SecureValidator) ValidateVolumeMount(volume VolumeMount) error {
	// Valida source
	if err := sv.pathValidator.ValidatePath(volume.Source); err != nil {
		return fmt.Errorf("volume source inválido: %w", err)
	}

	// Valida target
	if err := sv.pathValidator.ValidatePath(volume.Target); err != nil {
		return fmt.Errorf("volume target inválido: %w", err)
	}

	// Verifica se não está tentando montar diretórios sensíveis
	sensitivePaths := []string{
		"/etc", "/proc", "/sys", "/dev", "/boot", "/root",
		"/var/run/docker.sock", // Apenas se explicitamente necessário
	}

	for _, sensitive := range sensitivePaths {
		if strings.HasPrefix(volume.Target, sensitive) {
			return fmt.Errorf("mount em diretório sensível não permitido: %s", volume.Target)
		}
	}

	return nil
}

// ValidateResources valida limites de recursos
func (sv *SecureValidator) ValidateResources(resources Resources) error {
	// Valida CPU e Memory
	if err := security.ValidateResourceLimits(resources.CPUs, resources.Memory); err != nil {
		return err
	}

	// Valida GPU
	if resources.GPUs != "" && resources.GPUs != "all" {
		if err := security.ValidateResourceLimits(resources.GPUs, ""); err != nil {
			return fmt.Errorf("GPU limit inválido: %w", err)
		}
	}

	return nil
}

// ValidateSecret valida uma secret
func (sv *SecureValidator) ValidateSecret(secret Secret) error {
	// Valida nome
	cleanName, err := sv.inputSanitizer.SanitizeString(secret.Name)
	if err != nil {
		return fmt.Errorf("nome da secret inválido: %w", err)
	}

	if cleanName != secret.Name {
		return fmt.Errorf("nome da secret contém caracteres inválidos")
	}

	// Valida arquivo se especificado
	if secret.File != "" {
		if err := sv.pathValidator.ValidatePath(secret.File); err != nil {
			return fmt.Errorf("arquivo da secret inválido: %w", err)
		}
	}

	// Valida target
	if secret.Target != "" {
		if err := sv.pathValidator.ValidatePath(secret.Target); err != nil {
			return fmt.Errorf("target da secret inválido: %w", err)
		}
	}

	return nil
}

// ValidateRepositoryURL valida se uma URL de repositório é segura
func (sv *SecureValidator) ValidateRepositoryURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL do repositório não pode estar vazia")
	}

	// Limpar e sanitizar a URL
	cleanURL, err := sv.inputSanitizer.SanitizeString(url)
	if err != nil {
		return fmt.Errorf("erro ao sanitizar URL: %w", err)
	}
	if cleanURL != url {
		return fmt.Errorf("URL contém caracteres perigosos")
	}

	// Verificar se é uma URL válida do Git
	allowedPrefixes := []string{
		"https://github.com/",
		"https://gitlab.com/",
		"https://bitbucket.org/",
		"git@github.com:",
		"git@gitlab.com:",
		"git@bitbucket.org:",
	}

	valid := false
	for _, prefix := range allowedPrefixes {
		if strings.HasPrefix(url, prefix) {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("URL de repositório não permitida. Use GitHub, GitLab ou Bitbucket")
	}

	return nil
}
