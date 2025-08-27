package security

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// PathValidator valida caminhos contra path traversal
type PathValidator struct {
	allowedExtensions map[string]bool
	maxPathLength     int
}

// NewPathValidator cria um novo validador de paths
func NewPathValidator() *PathValidator {
	return &PathValidator{
		allowedExtensions: map[string]bool{
			".yml":        true,
			".yaml":       true,
			".env":        true,
			".dockerfile": true,
			".docker":     true,
		},
		maxPathLength: 255,
	}
}

// ValidatePath valida um caminho contra ataques de path traversal
func (pv *PathValidator) ValidatePath(path string) error {
	if path == "" {
		return fmt.Errorf("caminho vazio não é permitido")
	}

	// Verifica comprimento máximo
	if len(path) > pv.maxPathLength {
		return fmt.Errorf("caminho muito longo: máximo %d caracteres", pv.maxPathLength)
	}

	// Normaliza o caminho
	cleanPath := filepath.Clean(path)

	// Verifica se há tentativas de path traversal
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detectado: %s", path)
	}

	// Verifica caracteres perigosos
	if containsDangerousChars(cleanPath) {
		return fmt.Errorf("caracteres perigosos detectados no caminho: %s", path)
	}

	// Verifica se é um caminho absoluto suspeito
	if strings.HasPrefix(cleanPath, "/etc/") ||
		strings.HasPrefix(cleanPath, "/proc/") ||
		strings.HasPrefix(cleanPath, "/sys/") ||
		strings.HasPrefix(cleanPath, "/dev/") {
		return fmt.Errorf("acesso a diretório do sistema não permitido: %s", path)
	}

	return nil
}

// ValidateFileName valida um nome de arquivo
func (pv *PathValidator) ValidateFileName(filename string) error {
	if filename == "" {
		return fmt.Errorf("nome de arquivo vazio")
	}

	// Verifica extensão
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != "" && !pv.allowedExtensions[ext] {
		return fmt.Errorf("extensão de arquivo não permitida: %s", ext)
	}

	// Verifica caracteres perigosos no nome
	if containsDangerousChars(filename) {
		return fmt.Errorf("caracteres perigosos no nome do arquivo: %s", filename)
	}

	return nil
}

// containsDangerousChars verifica caracteres perigosos
func containsDangerousChars(input string) bool {
	dangerous := []string{
		"<", ">", ":", "\"", "|", "?", "*",
		"\x00", "\n", "\r", "\t",
	}

	for _, char := range dangerous {
		if strings.Contains(input, char) {
			return true
		}
	}

	// Verifica caracteres de controle
	for _, r := range input {
		if unicode.IsControl(r) && r != '\t' && r != '\n' && r != '\r' {
			return true
		}
	}

	return false
}

// InputSanitizer sanitiza inputs do usuário
type InputSanitizer struct {
	maxLength int
}

// NewInputSanitizer cria um novo sanitizador
func NewInputSanitizer(maxLength int) *InputSanitizer {
	return &InputSanitizer{
		maxLength: maxLength,
	}
}

// SanitizeString sanitiza uma string
func (is *InputSanitizer) SanitizeString(input string) (string, error) {
	if len(input) > is.maxLength {
		return "", fmt.Errorf("string muito longa: máximo %d caracteres", is.maxLength)
	}

	// Remove caracteres de controle
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\t' && r != '\n' && r != '\r' {
			return -1
		}
		return r
	}, input)

	// Trim espaços
	cleaned = strings.TrimSpace(cleaned)

	return cleaned, nil
}

// ValidateDomainName valida um nome de domínio
func ValidateDomainName(domain string) error {
	if domain == "" {
		return fmt.Errorf("domínio vazio")
	}

	if len(domain) > 253 {
		return fmt.Errorf("domínio muito longo")
	}

	// Regex para validar domínio
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]?(\.[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]?)*$`)
	if !domainRegex.MatchString(domain) {
		return fmt.Errorf("formato de domínio inválido: %s", domain)
	}

	return nil
}

// ValidateEmail valida um email
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email vazio")
	}

	if len(email) > 254 {
		return fmt.Errorf("email muito longo")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("formato de email inválido: %s", email)
	}

	return nil
}

// ValidateResourceLimits valida limites de recursos
func ValidateResourceLimits(cpus, memory string) error {
	if cpus != "" {
		cpuRegex := regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
		if !cpuRegex.MatchString(cpus) {
			return fmt.Errorf("formato de CPU inválido: %s", cpus)
		}
	}

	if memory != "" {
		memoryRegex := regexp.MustCompile(`^[0-9]+[kmg]?$`)
		if !memoryRegex.MatchString(strings.ToLower(memory)) {
			return fmt.Errorf("formato de memória inválido: %s", memory)
		}
	}

	return nil
}
