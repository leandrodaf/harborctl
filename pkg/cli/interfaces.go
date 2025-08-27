package cli

import "context"

// Command represents a CLI command
type Command interface {
	Name() string
	Description() string
	Execute(ctx context.Context, args []string) error
}

// Runner executa comandos
type Runner interface {
	Register(cmd Command)
	Run(ctx context.Context, args []string) error
}

// Output gerencia saídas
type Output interface {
	Info(msg string)
	Error(msg string)
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
