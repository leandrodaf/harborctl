package commands

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/leandrodaf/harborctl/internal/config"
	"github.com/leandrodaf/harborctl/pkg/cli"
)

// BeszelSetupCommand handles automatic Beszel setup via API
type BeszelSetupCommand struct {
	configManager config.Manager
	output        cli.Output
}

// NewBeszelSetupCommand creates a new Beszel setup command
func NewBeszelSetupCommand(configManager config.Manager, output cli.Output) cli.Command {
	return &BeszelSetupCommand{
		configManager: configManager,
		output:        output,
	}
}

func (c *BeszelSetupCommand) Name() string {
	return "beszel-setup"
}

func (c *BeszelSetupCommand) Description() string {
	return "🔧 Setup Beszel monitoring system with manual configuration guide"
}

func (c *BeszelSetupCommand) Usage() string {
	return `harborctl beszel-setup [options]

Guides you through Beszel monitoring system setup using Unix socket for optimal performance.
For same-host deployments (Hub and Agent in same compose), uses direct socket connection.

Options:
  -f, --file string      Configuration file path (default: stack.yml)
  --system-name string   Name for this system (default: hostname)
  -h, --help            Show help

Examples:
  harborctl beszel-setup
  harborctl beszel-setup --system-name production-server`
}

func (c *BeszelSetupCommand) Execute(ctx context.Context, args []string) error {
	var configFile, systemName string
	var help bool

	fs := flag.NewFlagSet("beszel-setup", flag.ContinueOnError)
	fs.StringVar(&configFile, "f", "stack.yml", "Configuration file path")
	fs.StringVar(&configFile, "file", "stack.yml", "Configuration file path")
	fs.StringVar(&systemName, "system-name", "", "Name for this system (default: hostname)")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.BoolVar(&help, "help", false, "Show help")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if help {
		c.output.Info(c.Usage())
		return nil
	}

	// Load existing configuration
	stack, err := c.configManager.Load(ctx, configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check if Beszel is enabled
	if !stack.Observability.Beszel.Enabled {
		c.output.Error("❌ Beszel is not enabled in this project")
		c.output.Info("   Enable Beszel first using: harborctl edit-server")
		return fmt.Errorf("beszel not enabled")
	}

	// Determine Hub URL
	var hubURL string
	if stack.Environment == "production" {
		hubURL = fmt.Sprintf("https://monitor.%s", stack.Domain)
	} else {
		hubURL = fmt.Sprintf("http://monitor.%s", stack.Domain)
	}

	// Default system name to hostname
	if systemName == "" {
		if hostname, err := os.Hostname(); err == nil {
			systemName = hostname
		} else {
			systemName = "harborctl-system"
		}
	}

	c.output.Info("🚀 Beszel Monitoring Setup Guide")
	c.output.Info("=" + strings.Repeat("=", 40))
	c.output.Info("")
	c.output.Infof("🌐 Hub URL: %s", hubURL)
	c.output.Infof("🖥️  System Name: %s", systemName)
	c.output.Info("")

	// Show setup instructions for Unix socket (same host)
	c.output.Info("📋 Unix Socket Setup (Optimal Performance):")
	c.output.Info("")
	c.output.Info("1. 🌐 Access Beszel Hub at:")
	c.output.Infof("   %s", hubURL)
	c.output.Info("")
	c.output.Info("2. 👤 Create your admin account (first visit)")
	c.output.Info("")
	c.output.Info("3. ➕ Add a new system with these settings:")
	c.output.Infof("   • Name: %s", systemName)
	c.output.Info("   • Host/IP: /beszel_socket/beszel.sock")
	c.output.Info("   • Port: (leave empty for Unix socket)")
	c.output.Info("")
	c.output.Info("4. ✅ Click 'Add System' - that's it!")
	c.output.Info("")
	c.output.Info("💡 No tokens needed for Unix socket connection!")
	c.output.Info("   The Hub connects directly to the Agent via shared socket.")
	c.output.Info("")
	c.output.Info("🚀 The system should appear online immediately.")

	return nil
}
