package docker

import "context"

type ComposeExecutor interface {
	ComposeUp(ctx context.Context, file string, build bool) error
	ComposeDown(ctx context.Context, file string) error
	ComposeStop(ctx context.Context, file string, timeout int) error
	ComposeStart(ctx context.Context, file string) error
	ComposeRestart(ctx context.Context, file string, timeout int) error
	ComposePause(ctx context.Context, file string) error
	ComposeUnpause(ctx context.Context, file string) error
}

type PruneExecutor interface {
	ImagePrune(ctx context.Context, filters ...string) error
	BuilderPrune(ctx context.Context, filters ...string) error
	VolumePrune(ctx context.Context) error
}

// Executor combines compose and prune operations
type Executor interface {
	ComposeExecutor
	PruneExecutor
}

type LifecycleManager interface {
	Deploy(ctx context.Context, composePath string, options DeployOptions) error
	Teardown(ctx context.Context, composePath string) error
	Stop(ctx context.Context, composePath string, timeout int) error
	Start(ctx context.Context, composePath string) error
	Restart(ctx context.Context, composePath string, timeout int) error
	Pause(ctx context.Context, composePath string) error
	Unpause(ctx context.Context, composePath string) error
}

type CleanupManager interface {
	Cleanup(ctx context.Context, options CleanupOptions) error
}

// Service combines lifecycle and cleanup operations
type Service interface {
	LifecycleManager
	CleanupManager
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
