package docker

import "context"

// Executor executa comandos Docker
type Executor interface {
	ComposeUp(ctx context.Context, file string, build bool) error
	ComposeDown(ctx context.Context, file string) error
	ImagePrune(ctx context.Context, filters ...string) error
	BuilderPrune(ctx context.Context, filters ...string) error
	VolumePrune(ctx context.Context) error
}

// Service representa operações Docker
type Service interface {
	Deploy(ctx context.Context, composePath string, options DeployOptions) error
	Teardown(ctx context.Context, composePath string) error
	Cleanup(ctx context.Context, options CleanupOptions) error
}

// DeployOptions configura o deploy
type DeployOptions struct {
	Build  bool
	Prune  bool
	Detach bool
}

// CleanupOptions configura a limpeza
type CleanupOptions struct {
	Images   bool
	Volumes  bool
	Networks bool
	MaxAge   string
}
