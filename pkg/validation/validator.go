package validation

import (
	"fmt"
	"strings"
)

type validator struct{}

// NewValidator creates a new validator
func NewValidator() Validator {
	return &validator{}
}

func (v *validator) ValidateHost(host string) error {
	if host == "" {
		return fmt.Errorf("host is required")
	}
	return nil
}

func (v *validator) ValidateAction(action string, validActions []string) error {
	for _, validAction := range validActions {
		if validAction == action {
			return nil
		}
	}
	return fmt.Errorf("invalid action: %s, valid actions: %s", action, strings.Join(validActions, ", "))
}

type actionValidator struct{}

// NewActionValidator creates a new action validator
func NewActionValidator() ActionValidator {
	return &actionValidator{}
}

func (av *actionValidator) IsValid(action string, validActions []string) bool {
	for _, validAction := range validActions {
		if validAction == action {
			return true
		}
	}
	return false
}

func (av *actionValidator) GetValidActions() []string {
	return []string{"status", "restart", "stop", "start", "details", "health"}
}
