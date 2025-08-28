package prompt

type TextPrompter interface {
	Text(message string, defaultValue ...string) (string, error)
	TextWithValidation(message string, validator ValidationFunc, defaultValue ...string) (string, error)
}

type PasswordPrompter interface {
	Password(message string) (string, error)
	PasswordWithValidation(message string, validator ValidationFunc) (string, error)
}

type ConfirmPrompter interface {
	Confirm(message string, defaultValue ...bool) (bool, error)
}

type SelectPrompter interface {
	Select(message string, options []string, defaultIndex ...int) (string, error)
	MultiSelect(message string, options []string) ([]string, error)
	InteractiveMultiSelect(message string, options []Option, preSelected ...int) ([]string, error)
}

type ValidatedPrompter interface {
	Email(message string, defaultValue ...string) (string, error)
	Domain(message string, defaultValue ...string) (string, error)
}

// Prompter combines all micro interfaces for convenience
type Prompter interface {
	TextPrompter
	PasswordPrompter
	ConfirmPrompter
	SelectPrompter
	ValidatedPrompter
}

// Option represents a selectable option
type Option struct {
	Label string
	Value string
}

// ValidationFunc validates user input
type ValidationFunc func(input string) error
