package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// executor implements Executor
type executor struct{}

// NewExecutor creates a new Docker executor
func NewExecutor() Executor {
	return &executor{}
}

func (e *executor) ComposeUp(ctx context.Context, file string, build bool) error {
	args := []string{"compose", "-f", file, "up", "-d"}
	if build {
		args = append(args, "--build")
	}
	return e.run(ctx, "docker", args...)
}

func (e *executor) ComposeDown(ctx context.Context, file string) error {
	return e.run(ctx, "docker", "compose", "-f", file, "down")
}

func (e *executor) ComposeStop(ctx context.Context, file string, timeout int) error {
	args := []string{"compose", "-f", file, "stop"}
	if timeout > 0 {
		args = append(args, "-t", fmt.Sprintf("%d", timeout))
	}
	return e.run(ctx, "docker", args...)
}

func (e *executor) ComposeStart(ctx context.Context, file string) error {
	return e.run(ctx, "docker", "compose", "-f", file, "start")
}

func (e *executor) ComposeRestart(ctx context.Context, file string, timeout int) error {
	args := []string{"compose", "-f", file, "restart"}
	if timeout > 0 {
		args = append(args, "-t", fmt.Sprintf("%d", timeout))
	}
	return e.run(ctx, "docker", args...)
}

func (e *executor) ComposePause(ctx context.Context, file string) error {
	return e.run(ctx, "docker", "compose", "-f", file, "pause")
}

func (e *executor) ComposeUnpause(ctx context.Context, file string) error {
	return e.run(ctx, "docker", "compose", "-f", file, "unpause")
}

func (e *executor) ImagePrune(ctx context.Context, filters ...string) error {
	args := []string{"image", "prune", "-af"}
	for _, filter := range filters {
		args = append(args, "--filter", filter)
	}
	return e.run(ctx, "docker", args...)
}

func (e *executor) BuilderPrune(ctx context.Context, filters ...string) error {
	args := []string{"builder", "prune", "-af"}
	for _, filter := range filters {
		args = append(args, "--filter", filter)
	}
	return e.run(ctx, "docker", args...)
}

func (e *executor) VolumePrune(ctx context.Context) error {
	return e.run(ctx, "docker", "volume", "prune", "-f")
}

func (e *executor) run(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// service implements Service
type service struct {
	executor Executor
}

// NewService creates a new Docker service
func NewService(executor Executor) Service {
	return &service{
		executor: executor,
	}
}

func (s *service) Deploy(ctx context.Context, composePath string, options DeployOptions) error {
	if err := s.executor.ComposeUp(ctx, composePath, options.Build); err != nil {
		return err
	}

	if options.Prune {
		return s.Cleanup(ctx, CleanupOptions{
			Images:  true,
			Volumes: true,
			MaxAge:  "168h",
		})
	}

	return nil
}

func (s *service) Teardown(ctx context.Context, composePath string) error {
	return s.executor.ComposeDown(ctx, composePath)
}

func (s *service) Stop(ctx context.Context, composePath string, timeout int) error {
	return s.executor.ComposeStop(ctx, composePath, timeout)
}

func (s *service) Start(ctx context.Context, composePath string) error {
	return s.executor.ComposeStart(ctx, composePath)
}

func (s *service) Restart(ctx context.Context, composePath string, timeout int) error {
	return s.executor.ComposeRestart(ctx, composePath, timeout)
}

func (s *service) Pause(ctx context.Context, composePath string) error {
	return s.executor.ComposePause(ctx, composePath)
}

func (s *service) Unpause(ctx context.Context, composePath string) error {
	return s.executor.ComposeUnpause(ctx, composePath)
}

func (s *service) Cleanup(ctx context.Context, options CleanupOptions) error {
	if options.Images {
		if err := s.executor.ImagePrune(ctx, "until="+options.MaxAge); err != nil {
			return err
		}
		if err := s.executor.BuilderPrune(ctx, "until="+options.MaxAge); err != nil {
			return err
		}
	}

	if options.Volumes {
		if err := s.executor.VolumePrune(ctx); err != nil {
			return err
		}
	}

	return nil
}
