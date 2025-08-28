package main

import (
	"context"
	"os"

	"github.com/leandrodaf/harborctl/internal/commands"
	"github.com/leandrodaf/harborctl/internal/compose"
	"github.com/leandrodaf/harborctl/internal/config"
	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/docker"
	"github.com/leandrodaf/harborctl/pkg/fs"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	ctx := context.Background()

	// Initialize output
	output := cli.NewOutput()

	// Handle version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		output.Infof("harborctl version %s", version)
		output.Infof("Built: %s", buildTime)
		output.Infof("Commit: %s", gitCommit)
		return
	}

	// Handle help flag
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		showHelp(output)
		return
	}

	// Initialize dependencies with proper error handling
	dependencies, err := initializeDependencies()
	if err != nil {
		output.Errorf("Failed to initialize dependencies: %v", err)
		os.Exit(1)
	}

	// Initialize CLI runner
	runner := cli.NewRunner(output)

	// Register all commands
	registerCommands(runner, dependencies.ConfigManager, dependencies.ComposeService, dependencies.DockerService, dependencies.FileSystem, output)

	// Run the CLI
	if err := runner.Run(ctx, os.Args[1:]); err != nil {
		output.Errorf("Error: %v", err)
		os.Exit(1)
	}
}

// Dependencies holds all service dependencies
type Dependencies struct {
	ConfigManager  config.Manager
	ComposeService compose.Service
	DockerService  docker.Service
	FileSystem     fs.FileSystem
}

func initializeDependencies() (*Dependencies, error) {
	// Initialize filesystem
	filesystem := fs.NewFileSystem()

	// Initialize config components
	configLoader := fs.NewConfigLoader(filesystem)
	validator := config.NewValidator()
	configManager := config.NewManager(configLoader, filesystem, validator)

	// Initialize docker components
	dockerExecutor := docker.NewExecutor()
	dockerService := docker.NewService(dockerExecutor)

	// Initialize compose service
	composeService := compose.NewDefaultService()

	return &Dependencies{
		ConfigManager:  configManager,
		ComposeService: composeService,
		DockerService:  dockerService,
		FileSystem:     filesystem,
	}, nil
}

func registerCommands(
	runner cli.Runner,
	configManager config.Manager,
	composeService compose.Service,
	dockerService docker.Service,
	filesystem fs.FileSystem,
	output cli.Output,
) {
	// Register init command
	runner.Register(commands.NewInitCommand(configManager, output))

	// Register init-server command
	runner.Register(commands.NewInitServerCommand(configManager, output))

	// Register deploy-service command
	runner.Register(commands.NewDeployServiceCommand(configManager, composeService, dockerService, filesystem, output))

	// Register up command
	runner.Register(commands.NewUpCommand(configManager, composeService, dockerService, filesystem, output))

	// Register down command
	runner.Register(commands.NewDownCommand(dockerService, output))

	// Register stop command
	runner.Register(commands.NewStopCommand(dockerService, output))

	// Register start command
	runner.Register(commands.NewStartCommand(dockerService, output))

	// Register restart command
	runner.Register(commands.NewRestartCommand(dockerService, output))

	// Register pause command
	runner.Register(commands.NewPauseCommand(dockerService, output))

	// Register unpause command
	runner.Register(commands.NewUnpauseCommand(dockerService, output))

	// Register status command
	runner.Register(commands.NewStatusCommand(dockerService, output))

	// Register logs command
	runner.Register(commands.NewLogsCommand(dockerService, output))

	// Register remote-logs command
	runner.Register(commands.NewRemoteLogsCommand(output))

	// Register remote-control command
	runner.Register(commands.NewRemoteControlCommand(output))

	// Register scale command
	runner.Register(commands.NewScaleCommand(configManager, dockerService, output))

	// Register validate command
	runner.Register(commands.NewValidateCommand(configManager, output))

	// Register render command
	runner.Register(commands.NewRenderCommand(configManager, composeService, filesystem, output))

	// Register hash-password command
	runner.Register(commands.NewHashPasswordCommand(output))

	// Register security-audit command
	runner.Register(commands.NewSecurityAuditCommand(configManager, output))

	// Register docs command
	runner.Register(commands.NewDocsCommand(output))
}

func showHelp(output cli.Output) {
	output.Info("🚢 Harbor CLI - Deployment Tool")
	output.Info("")
	output.Info("USAGE:")
	output.Info("  harborctl [command] [flags]")
	output.Info("")
	output.Info("COMMANDS:")
	output.Info("  init              Initialize new project configuration")
	output.Info("  init-server       Initialize server with required dependencies")
	output.Info("  deploy-service    Deploy a service to the stack")
	output.Info("")
	output.Info("LIFECYCLE:")
	output.Info("  up                Start all services in the stack")
	output.Info("  down              Stop and remove all services")
	output.Info("  stop              Stop services without removing containers")
	output.Info("  start             Start previously stopped services")
	output.Info("  restart           Restart all services")
	output.Info("  pause             Pause all services")
	output.Info("  unpause           Unpause all services")
	output.Info("")
	output.Info("MANAGEMENT:")
	output.Info("  status            Show status of all services")
	output.Info("  logs              Show services logs")
	output.Info("  scale             Scale services up or down")
	output.Info("")
	output.Info("REMOTE:")
	output.Info("  remote-logs       View logs from remote server")
	output.Info("  remote-control    Control services on remote server")
	output.Info("")
	output.Info("TOOLS:")
	output.Info("  validate          Validate stack configuration")
	output.Info("  render            Render docker-compose configuration")
	output.Info("  hash-password     Generate hashed password for authentication")
	output.Info("  security-audit    Run security audit on the stack")
	output.Info("  docs              Show documentation and guides")
	output.Info("")
	output.Info("FLAGS:")
	output.Info("  -h, --help        Show this help message")
	output.Info("  -v, --version     Show version information")
	output.Info("")
	output.Info("For more information about a specific command, run:")
	output.Info("  harborctl [command] --help")
}
