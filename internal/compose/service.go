package compose

import (
	"context"

	"github.com/leandrodaf/harborctl/internal/config"
)

// Generator gera docker-compose
type Generator interface {
	Generate(ctx context.Context, stack *config.Stack, options GenerateOptions) ([]byte, error)
}

// Service gerencia geração de compose
type Service interface {
	Generate(ctx context.Context, stack *config.Stack, options GenerateOptions) ([]byte, error)
}

// GenerateOptions configura a geração
type GenerateOptions struct {
	DisableDozzle bool
	DisableBeszel bool
}

// service implementa Service
type service struct {
	generator Generator
}

// NewService cria um novo serviço de compose
func NewService(generator Generator) Service {
	return &service{
		generator: generator,
	}
}

func (s *service) Generate(ctx context.Context, stack *config.Stack, options GenerateOptions) ([]byte, error) {
	return s.generator.Generate(ctx, stack, options)
}
