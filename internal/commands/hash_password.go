package commands

import (
	"context"
	"flag"
	"fmt"

	"github.com/leandrodaf/harborctl/pkg/auth"
	"github.com/leandrodaf/harborctl/pkg/cli"
)

// hashPasswordCommand implementa o comando hash-password
type hashPasswordCommand struct {
	output cli.Output
}

// NewHashPasswordCommand cria um novo comando hash-password
func NewHashPasswordCommand(output cli.Output) cli.Command {
	return &hashPasswordCommand{
		output: output,
	}
}

func (c *hashPasswordCommand) Name() string {
	return "hash-password"
}

func (c *hashPasswordCommand) Description() string {
	return "Gera hash bcrypt para senha (para basic auth)"
}

func (c *hashPasswordCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("hash-password", flag.ExitOnError)

	var password string
	var generate bool
	var length int

	fs.StringVar(&password, "password", "", "senha para gerar hash")
	fs.BoolVar(&generate, "generate", false, "gerar senha aleatória")
	fs.IntVar(&length, "length", 12, "tamanho da senha gerada")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if generate {
		generatedPassword, err := auth.GenerateRandomPassword(length)
		if err != nil {
			return fmt.Errorf("erro ao gerar senha: %w", err)
		}
		password = generatedPassword
		c.output.Infof("Senha gerada: %s", password)
	}

	if password == "" {
		return fmt.Errorf("especifique uma senha com --password ou use --generate")
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("erro ao gerar hash: %w", err)
	}

	c.output.Info("Hash bcrypt:")
	c.output.Info(hash)
	c.output.Info("")
	c.output.Info("Exemplo de uso no stack.yml:")
	c.output.Infof(`  basic_auth:
    enabled: true
    users:
      admin: "%s"`, hash)

	return nil
}
