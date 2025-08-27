package middleware

import (
	"context"
	"time"

	"github.com/leandrodaf/harborctl/pkg/cli"
)

// LoggingMiddleware adds logging to commands
type LoggingMiddleware struct {
	output cli.Output
	next   cli.Command
}

// NewLoggingMiddleware creates a logging middleware
func NewLoggingMiddleware(next cli.Command, output cli.Output) cli.Command {
	return &LoggingMiddleware{
		output: output,
		next:   next,
	}
}

func (m *LoggingMiddleware) Name() string {
	return m.next.Name()
}

func (m *LoggingMiddleware) Description() string {
	return m.next.Description()
}

func (m *LoggingMiddleware) Execute(ctx context.Context, args []string) error {
	start := time.Now()
	m.output.Infof("Starting command: %s", m.next.Name())

	err := m.next.Execute(ctx, args)

	duration := time.Since(start)
	if err != nil {
		m.output.Errorf("Command %s failed after %v: %v", m.next.Name(), duration, err)
	} else {
		m.output.Infof("Command %s completed in %v", m.next.Name(), duration)
	}

	return err
}

// TimingMiddleware adds time measurement
type TimingMiddleware struct {
	output cli.Output
	next   cli.Command
}

// NewTimingMiddleware creates a timing middleware
func NewTimingMiddleware(next cli.Command, output cli.Output) cli.Command {
	return &TimingMiddleware{
		output: output,
		next:   next,
	}
}

func (m *TimingMiddleware) Name() string {
	return m.next.Name()
}

func (m *TimingMiddleware) Description() string {
	return m.next.Description()
}

func (m *TimingMiddleware) Execute(ctx context.Context, args []string) error {
	start := time.Now()

	err := m.next.Execute(ctx, args)

	duration := time.Since(start)
	m.output.Infof("⏱️  Tempo de execução: %v", duration)

	return err
}
