package commands

import (
	"context"
	"flag"
	"fmt"

	"github.com/leandrodaf/harborctl/internal/config"
	"github.com/leandrodaf/harborctl/internal/crypto"
	"github.com/leandrodaf/harborctl/pkg/cli"
)

// RegenerateBeszelKeysCommand handles regenerating Beszel keys for existing projects
type RegenerateBeszelKeysCommand struct {
	configManager config.Manager
	output        cli.Output
}

// NewRegenerateBeszelKeysCommand creates a new regenerate Beszel keys command
func NewRegenerateBeszelKeysCommand(configManager config.Manager, output cli.Output) cli.Command {
	return &RegenerateBeszelKeysCommand{
		configManager: configManager,
		output:        output,
	}
}

func (c *RegenerateBeszelKeysCommand) Name() string {
	return "regenerate-beszel-keys"
}

func (c *RegenerateBeszelKeysCommand) Description() string {
	return "Regenerate Beszel authentication keys for existing project"
}

func (c *RegenerateBeszelKeysCommand) Usage() string {
	return `harborctl regenerate-beszel-keys [options]

Regenerates Beszel SSH keys and token for an existing project.

Options:
  -f, --file string    Configuration file path (default: stack.yml)
  -h, --help          Show help

Examples:
  harborctl regenerate-beszel-keys
  harborctl regenerate-beszel-keys -f custom-stack.yml`
}

func (c *RegenerateBeszelKeysCommand) Execute(ctx context.Context, args []string) error {
	var configFile string
	var help bool

	fs := flag.NewFlagSet("regenerate-beszel-keys", flag.ContinueOnError)
	fs.StringVar(&configFile, "f", "stack.yml", "Configuration file path")
	fs.StringVar(&configFile, "file", "stack.yml", "Configuration file path")
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
		c.output.Info("❌ Beszel is not enabled in this project")
		c.output.Info("   Enable Beszel first using: harborctl edit-server")
		return nil
	}

	c.output.Info("🔐 Regenerating Beszel authentication keys...")

	// Generate new SSH key pair
	pubKey, _, err := crypto.GenerateED25519KeyPair()
	if err != nil {
		return fmt.Errorf("failed to generate SSH keys: %w", err)
	}

	// Generate new token
	token, err := crypto.GenerateBeszelToken()
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Update configuration
	stack.Observability.Beszel.PublicKey = pubKey
	stack.Observability.Beszel.Token = token

	// Save updated configuration
	if err := c.configManager.SaveBaseConfig(ctx, configFile, stack); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	c.output.Info("✅ Beszel keys regenerated successfully!")
	c.output.Info("")
	c.output.Info("🔑 New SSH Key: " + pubKey[:50] + "...")
	c.output.Info("🎫 New Token: " + token[:20] + "...")
	c.output.Info("")
	c.output.Info("🚀 Next steps:")
	c.output.Info("   1. Run: harborctl up")
	c.output.Info("   2. The agent should now connect successfully")

	return nil
}
