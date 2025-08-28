package prompt

import (
	"fmt"
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

// ValidateDomain validates domain format
func ValidateDomain(input string) error {
	if input == "" {
		return fmt.Errorf("domain is required")
	}

	// Basic domain validation
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	if !domainRegex.MatchString(input) {
		return fmt.Errorf("invalid domain format")
	}

	// Check for valid TLD
	parts := strings.Split(input, ".")
	if len(parts) < 2 {
		return fmt.Errorf("domain must have at least two parts")
	}

	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(input string) error {
	if len(input) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	// Check for at least one number and one letter
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(input)
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(input)

	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}

	if !hasLetter {
		return fmt.Errorf("password must contain at least one letter")
	}

	return nil
}

// ValidateProjectName validates project name format
func ValidateProjectName(input string) error {
	if err := ValidateRequired(input); err != nil {
		return err
	}

	// Project names should be lowercase alphanumeric with hyphens
	projectRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !projectRegex.MatchString(input) {
		return fmt.Errorf("project name must contain only lowercase letters, numbers, and hyphens")
	}

	if strings.HasPrefix(input, "-") || strings.HasSuffix(input, "-") {
		return fmt.Errorf("project name cannot start or end with a hyphen")
	}

	return nil
}

// ValidateSubdomain validates subdomain format
func ValidateSubdomain(input string) error {
	if input == "" {
		return fmt.Errorf("subdomain is required")
	}

	// Subdomain validation (simpler than full domain)
	subdomainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`)
	if !subdomainRegex.MatchString(input) {
		return fmt.Errorf("invalid subdomain format")
	}

	return nil
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
