package cli

import (
	"context"
	"fmt"
	"os"
)

// runner implementa Runner
type runner struct {
	commands map[string]Command
	output   Output
}

// NewRunner cria um novo runner
func NewRunner(output Output) Runner {
	return &runner{
		commands: make(map[string]Command),
		output:   output,
	}
}

// Register registra um comando
func (r *runner) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

// Run executa um comando
func (r *runner) Run(ctx context.Context, args []string) error {
	if len(args) < 1 {
		r.showUsage()
		return fmt.Errorf("comando não especificado")
	}

	cmdName := args[0]
	cmd, exists := r.commands[cmdName]
	if !exists {
		r.showUsage()
		return fmt.Errorf("comando desconhecido: %s", cmdName)
	}

	return cmd.Execute(ctx, args[1:])
}

func (r *runner) showUsage() {
	r.output.Info("harborctl - Gerenciador de Stack Docker")
	r.output.Info("Comandos disponíveis:")
	for name, cmd := range r.commands {
		r.output.Infof("  %-10s %s", name, cmd.Description())
	}
}

// output implementa Output
type output struct{}

// NewOutput cria um novo output
func NewOutput() Output {
	return &output{}
}

func (o *output) Info(msg string) {
	fmt.Println(msg)
}

func (o *output) Error(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

func (o *output) Infof(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func (o *output) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
