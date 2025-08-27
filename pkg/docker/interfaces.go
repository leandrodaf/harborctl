package docker

import "context"

// Executor executes Docker commands
type Executor interface {
	ComposeUp(ctx context.Context, file string, build bool) error
	ComposeDown(ctx context.Context, file string) error
	ImagePrune(ctx context.Context, filters ...string) error
	BuilderPrune(ctx context.Context, filters ...string) error
	VolumePrune(ctx context.Context) error
}

// Service represents Docker operations
type Service interface {
	Deploy(ctx context.Context, composePath string, options DeployOptions) error
	Teardown(ctx context.Context, composePath string) error
	Cleanup(ctx context.Context, options CleanupOptions) error
}

// DeployOptions configures deployment
type DeployOptions struct {
	Build  bool
	Prune  bool
	Detach bool
}

// CleanupOptions configures cleanup
type CleanupOptions struct {
	Images   bool
	Volumes  bool
	Networks bool
	MaxAge   string
}
