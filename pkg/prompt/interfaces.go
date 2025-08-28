package prompt

// Prompter handles interactive user prompts
type Prompter interface {
	// Text prompts for text input with optional default
	Text(message string, defaultValue ...string) (string, error)

	// TextWithValidation prompts for text input with validation
	TextWithValidation(message string, validator ValidationFunc, defaultValue ...string) (string, error)

	// Password prompts for password input (hidden)
	Password(message string) (string, error)

	// PasswordWithValidation prompts for password input with validation
	PasswordWithValidation(message string, validator ValidationFunc) (string, error)

	// Confirm prompts for yes/no confirmation
	Confirm(message string, defaultValue ...bool) (bool, error)

	// Select prompts for single selection from options
	Select(message string, options []string, defaultIndex ...int) (string, error)

	// MultiSelect prompts for multiple selections from options
	MultiSelect(message string, options []string) ([]string, error)

	// InteractiveMultiSelect provides an interactive multi-selection interface
	InteractiveMultiSelect(message string, options []Option, preSelected ...int) ([]string, error)

	// Email prompts for email input with built-in validation
	Email(message string, defaultValue ...string) (string, error)

	// Domain prompts for domain input with built-in validation
	Domain(message string, defaultValue ...string) (string, error)
}

// Option represents a selectable option
type Option struct {
	Label string
	Value string
}

// ValidationFunc validates user input
type ValidationFunc func(input string) error
