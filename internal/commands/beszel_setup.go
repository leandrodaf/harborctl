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

Guides you through Beszel monitoring system setup by:
1. Displaying connection information for manual configuration
2. Providing step-by-step instructions for web interface
3. Optionally updating stack configuration with new tokens

Options:
  -f, --file string      Configuration file path (default: stack.yml)
  --token string         Agent token from Beszel Hub (optional)
  --public-key string    Public key from Beszel Hub (optional)
  --system-name string   Name for this system (default: hostname)
  -h, --help            Show help

Examples:
  harborctl beszel-setup
  harborctl beszel-setup --token <token> --public-key <key>
  harborctl beszel-setup --system-name production-server`
}

func (c *BeszelSetupCommand) Execute(ctx context.Context, args []string) error {
	var configFile, token, publicKey, systemName string
	var help bool

	fs := flag.NewFlagSet("beszel-setup", flag.ContinueOnError)
	fs.StringVar(&configFile, "f", "stack.yml", "Configuration file path")
	fs.StringVar(&configFile, "file", "stack.yml", "Configuration file path")
	fs.StringVar(&token, "token", "", "Agent token from Beszel Hub")
	fs.StringVar(&publicKey, "public-key", "", "Public key from Beszel Hub")
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
	c.output.Info("=" + strings.Repeat("=", 35))
	c.output.Info("")
	c.output.Infof("🌐 Hub URL: %s", hubURL)
	c.output.Infof("🖥️  System Name: %s", systemName)
	c.output.Info("")

	if token != "" && publicKey != "" {
		// Update configuration with provided tokens
		stack.Observability.Beszel.Token = token
		stack.Observability.Beszel.PublicKey = publicKey

		if err := c.configManager.SaveBaseConfig(ctx, configFile, stack); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		c.output.Info("✅ Configuration updated successfully!")
		c.output.Info("")
		c.output.Info("🚀 Next steps:")
		c.output.Info("   1. Run: harborctl up")
		c.output.Info("   2. Your system should now appear in Beszel Hub")

		return nil
	}

	// Show manual setup instructions
	c.output.Info("📋 Manual Setup Instructions:")
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
	c.output.Info("4. 📋 Copy the generated token and public key")
	c.output.Info("")
	c.output.Info("5. 🔧 Run this command with your tokens:")
	c.output.Info("   harborctl beszel-setup \\")
	c.output.Info("     --token <your-token> \\")
	c.output.Info("     --public-key <your-public-key>")
	c.output.Info("")
	c.output.Info("6. 🚀 Restart services:")
	c.output.Info("   harborctl up")
	c.output.Info("")
	c.output.Info("💡 The system will connect automatically via Unix socket")
	c.output.Info("   for optimal performance and security.")

	return nil
}
