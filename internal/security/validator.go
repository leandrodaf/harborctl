package security

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// PathValidator validates paths against path traversal
type PathValidator struct {
	allowedExtensions map[string]bool
	maxPathLength     int
}

// NewPathValidator creates a new path validator
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

// ValidatePath validates a path against path traversal attacks
func (pv *PathValidator) ValidatePath(path string) error {
	if path == "" {
		return fmt.Errorf("empty path not allowed")
	}

	// Check maximum length
	if len(path) > pv.maxPathLength {
		return fmt.Errorf("path too long: maximum %d characters", pv.maxPathLength)
	}

	// Normalize path
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected: %s", path)
	}

	// Check dangerous characters
	if containsDangerousChars(cleanPath) {
		return fmt.Errorf("dangerous characters detected in path: %s", path)
	}

	// Check if it's a suspicious absolute path
	if strings.HasPrefix(cleanPath, "/etc/") ||
		strings.HasPrefix(cleanPath, "/proc/") ||
		strings.HasPrefix(cleanPath, "/sys/") ||
		strings.HasPrefix(cleanPath, "/dev/") {
		return fmt.Errorf("access to system directory not allowed: %s", path)
	}

	return nil
}

// ValidateFileName validates a filename
func (pv *PathValidator) ValidateFileName(filename string) error {
	if filename == "" {
		return fmt.Errorf("empty filename")
	}

	// Check extension
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != "" && !pv.allowedExtensions[ext] {
		return fmt.Errorf("file extension not allowed: %s", ext)
	}

	// Check dangerous characters in name
	if containsDangerousChars(filename) {
		return fmt.Errorf("dangerous characters in filename: %s", filename)
	}

	return nil
}

// containsDangerousChars checks for dangerous characters
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

	// Check for control characters
	for _, r := range input {
		if unicode.IsControl(r) && r != '\t' && r != '\n' && r != '\r' {
			return true
		}
	}

	return false
}

// InputSanitizer sanitizes user inputs
type InputSanitizer struct {
	maxLength int
}

// NewInputSanitizer creates a new sanitizer
func NewInputSanitizer(maxLength int) *InputSanitizer {
	return &InputSanitizer{
		maxLength: maxLength,
	}
}

// SanitizeString sanitizes a string
func (is *InputSanitizer) SanitizeString(input string) (string, error) {
	if len(input) > is.maxLength {
		return "", fmt.Errorf("string too long: maximum %d characters", is.maxLength)
	}

	// Remove control characters
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\t' && r != '\n' && r != '\r' {
			return -1
		}
		return r
	}, input)

	// Trim spaces
	cleaned = strings.TrimSpace(cleaned)

	return cleaned, nil
}

// ValidateDomainName validates a domain name
func ValidateDomainName(domain string) error {
	if domain == "" {
		return fmt.Errorf("empty domain")
	}

	if len(domain) > 253 {
		return fmt.Errorf("domain too long")
	}

	// Regex to validate domain
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]?(\.[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]?)*$`)
	if !domainRegex.MatchString(domain) {
		return fmt.Errorf("invalid domain format: %s", domain)
	}

	return nil
}

// ValidateEmail validates an email
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("empty email")
	}

	if len(email) > 254 {
		return fmt.Errorf("email too long")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}

	return nil
}

// ValidateResourceLimits validates resource limits
func ValidateResourceLimits(cpus, memory string) error {
	if cpus != "" {
		cpuRegex := regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
		if !cpuRegex.MatchString(cpus) {
			return fmt.Errorf("invalid CPU format: %s", cpus)
		}
	}

	if memory != "" {
		memoryRegex := regexp.MustCompile(`^[0-9]+[kmg]?$`)
		if !memoryRegex.MatchString(strings.ToLower(memory)) {
			return fmt.Errorf("invalid memory format: %s", memory)
		}
	}

	return nil
}
