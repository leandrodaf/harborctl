package prompt

import (
	"fmt"
	"net"
	"net/mail"
	"regexp"
	"strings"
)

// Common validation functions

// ValidateRequired validates that input is not empty
func ValidateRequired(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("required value")
	}
	return nil
}

// ValidateEmail validates email format
func ValidateEmail(input string) error {
	if _, err := mail.ParseAddress(input); err != nil {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidateDomain validates domain name format
func ValidateDomain(input string) error {
	if input == "" {
		return fmt.Errorf("domain is required")
	}

	// Remove protocol if present
	input = strings.TrimPrefix(input, "http://")
	input = strings.TrimPrefix(input, "https://")

	// Remove trailing slash
	input = strings.TrimSuffix(input, "/")

	// Basic domain regex
	domainRegex := `^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`
	matched, err := regexp.MatchString(domainRegex, input)
	if err != nil {
		return fmt.Errorf("domain validation error")
	}
	if !matched {
		return fmt.Errorf("invalid domain format")
	}

	// Check if it's a valid hostname (additional validation)
	if net.ParseIP(input) == nil && !strings.Contains(input, ".") && input != "localhost" {
		return fmt.Errorf("domain must contain at least one dot (.) or be 'localhost'")
	}

	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(input string) error {
	if len(input) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	return nil
}

// ValidatePort validates port number
func ValidatePort(input string) error {
	if input == "" {
		return fmt.Errorf("port is required")
	}

	// Basic range check for common ports
	if len(input) > 5 {
		return fmt.Errorf("invalid port")
	}

	// Check if it's numeric (simplified check)
	for _, char := range input {
		if char < '0' || char > '9' {
			return fmt.Errorf("port must contain only numbers")
		}
	}

	return nil
}

// ValidateSubdomain validates subdomain format
func ValidateSubdomain(input string) error {
	if input == "" {
		return fmt.Errorf("subdomain is required")
	}

	// Subdomain regex (letters, numbers, hyphens, no dots)
	subdomainRegex := `^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`
	matched, err := regexp.MatchString(subdomainRegex, input)
	if err != nil {
		return fmt.Errorf("subdomain validation error")
	}
	if !matched {
		return fmt.Errorf("invalid subdomain. Use only letters, numbers and hyphens")
	}

	return nil
}

// ValidateProjectName validates project name format
func ValidateProjectName(input string) error {
	if input == "" {
		return fmt.Errorf("project name is required")
	}

	// Project name regex (letters, numbers, hyphens, underscores)
	projectRegex := `^[a-zA-Z0-9][a-zA-Z0-9_\-]{0,49}$`
	matched, err := regexp.MatchString(projectRegex, input)
	if err != nil {
		return fmt.Errorf("project name validation error")
	}
	if !matched {
		return fmt.Errorf("invalid project name. Use letters, numbers, _ and - (max 50 chars)")
	}

	return nil
}

// ValidateURL validates URL format
func ValidateURL(input string) error {
	if input == "" {
		return fmt.Errorf("URL is required")
	}

	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}

	return nil
}

// ValidateIPAddress validates IP address format (basic)
func ValidateIPAddress(input string) error {
	if input == "" {
		return fmt.Errorf("IP address is required")
	}

	parts := strings.Split(input, ".")
	if len(parts) != 4 {
		return fmt.Errorf("IP address must have 4 octets separated by dots")
	}

	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return fmt.Errorf("invalid octet in IP")
		}
		for _, char := range part {
			if char < '0' || char > '9' {
				return fmt.Errorf("IP must contain only numbers and dots")
			}
		}
	}

	return nil
}

// ValidateNotEmpty validates that input is not empty
func ValidateNotEmpty(input string) error {
	return ValidateRequired(input)
}

// ValidateMinLength validates minimum length
func ValidateMinLength(minLen int) ValidationFunc {
	return func(input string) error {
		if len(input) < minLen {
			return fmt.Errorf("must be at least %d characters long", minLen)
		}
		return nil
	}
}

// ValidateMaxLength validates maximum length
func ValidateMaxLength(maxLen int) ValidationFunc {
	return func(input string) error {
		if len(input) > maxLen {
			return fmt.Errorf("must be at most %d characters long", maxLen)
		}
		return nil
	}
}

// ValidateAlphanumeric validates alphanumeric input
func ValidateAlphanumeric(input string) error {
	if input == "" {
		return fmt.Errorf("value is required")
	}

	for _, char := range input {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return fmt.Errorf("must contain only letters and numbers")
		}
	}

	return nil
}

// ValidateOneOf validates that input is one of the allowed values
func ValidateOneOf(allowedValues ...string) ValidationFunc {
	return func(input string) error {
		for _, allowed := range allowedValues {
			if input == allowed {
				return nil
			}
		}
		return fmt.Errorf("must be one of: %s", strings.Join(allowedValues, ", "))
	}
}

// ValidateRegex validates input against a regular expression
func ValidateRegex(pattern string, errorMsg string) ValidationFunc {
	return func(input string) error {
		matched, err := regexp.MatchString(pattern, input)
		if err != nil {
			return fmt.Errorf("validation error: %v", err)
		}
		if !matched {
			return fmt.Errorf("%s", errorMsg)
		}
		return nil
	}
}

// CombineValidators combines multiple validators
func CombineValidators(validators ...ValidationFunc) ValidationFunc {
	return func(input string) error {
		for _, validator := range validators {
			if err := validator(input); err != nil {
				return err
			}
		}
		return nil
	}
}
