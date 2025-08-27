package middleware

import (
	"context"
	"time"

	"github.com/leandrodaf/harborctl/pkg/cli"
)

// LoggingMiddleware adiciona logging aos comandos
type LoggingMiddleware struct {
	output cli.Output
	next   cli.Command
}

// NewLoggingMiddleware cria um middleware de logging
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
	m.output.Infof("Iniciando comando: %s", m.next.Name())

	err := m.next.Execute(ctx, args)

	duration := time.Since(start)
	if err != nil {
		m.output.Errorf("Comando %s falhou após %v: %v", m.next.Name(), duration, err)
	} else {
		m.output.Infof("Comando %s concluído em %v", m.next.Name(), duration)
	}

	return err
}

// TimingMiddleware adiciona medição de tempo
type TimingMiddleware struct {
	output cli.Output
	next   cli.Command
}

// NewTimingMiddleware cria um middleware de timing
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
