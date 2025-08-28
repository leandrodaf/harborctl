package commands

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/leandrodaf/harborctl/internal/config"
	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/prompt"
	"golang.org/x/crypto/bcrypt"
)

// SetupCommand handles interactive server setup
type SetupCommand struct {
	configManager config.Manager
	prompter      prompt.Prompter
	errorHandler  *prompt.ErrorHandler
	output        cli.Output
}

// NewSetupCommand creates a new interactive setup command
func NewSetupCommand(configManager config.Manager, output cli.Output) cli.Command {
	prompter := prompt.NewPrompter()
	return &SetupCommand{
		configManager: configManager,
		prompter:      prompter,
		errorHandler:  prompt.NewErrorHandler(prompter),
		output:        output,
	}
}

func (c *SetupCommand) Name() string {
	return "setup"
}

func (c *SetupCommand) Description() string {
	return "Interactive server setup wizard"
}

func (c *SetupCommand) Execute(ctx context.Context, args []string) error {
	// Add panic recovery
	defer c.errorHandler.RecoverFromPanic()

	return c.errorHandler.SafeOperation(ctx, "Interactive Setup", func() error {
		return c.runSetup(ctx)
	})
}

func (c *SetupCommand) runSetup(ctx context.Context) error {
	c.output.Info("🚀 Welcome to HarborCtl Interactive Setup!")
	c.output.Info("This wizard will help you configure your server infrastructure.")
	c.output.Info("")

	// Basic Configuration
	domain, err := c.prompter.Domain("Enter your domain (e.g., example.com)")
	if err != nil {
		return fmt.Errorf("failed to get domain: %w", err)
	}

	email, err := c.prompter.Email("Enter your email for SSL certificates", fmt.Sprintf("admin@%s", domain))
	if err != nil {
		return fmt.Errorf("failed to get email: %w", err)
	}

	project, err := c.prompter.TextWithValidation("Enter project name", prompt.ValidateProjectName, "deploy")
	if err != nil {
		return fmt.Errorf("failed to get project name: %w", err)
	}

	// SSL Configuration
	sslMode, err := c.prompter.Select("Choose SSL mode", []string{
		"Automatic SSL (Let's Encrypt) - Recommended for production",
		"Disabled - For local development only",
	}, 0)
	if err != nil {
		return fmt.Errorf("failed to get SSL mode: %w", err)
	}

	// Environment Detection
	env := "production"
	if strings.Contains(sslMode, "Disabled") {
		env = "local"
	}

	// Observability Configuration
	enableObservability, err := c.prompter.Confirm("Enable observability services (logs + monitoring)?", true)
	if err != nil {
		return fmt.Errorf("failed to get observability preference: %w", err)
	}

	var observability config.Observability
	if enableObservability {
		observability, err = c.configureObservability()
		if err != nil {
			return fmt.Errorf("failed to configure observability: %w", err)
		}
	}

	// Build configuration
	stack := c.buildStack(domain, email, project, env, observability)

	// Save configuration
	filename := "server-base.yml"
	if err := c.configManager.SaveBaseConfig(ctx, filename, stack); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	// Show summary
	c.showSummary(stack, filename)

	return nil
}

func (c *SetupCommand) configureObservability() (config.Observability, error) {
	var obs config.Observability

	// Dozzle Configuration
	enableDozzle, err := c.prompter.Confirm("Enable Dozzle (log viewer)?", true)
	if err != nil {
		return obs, err
	}

	if enableDozzle {
		obs.Dozzle.Enabled = true
		obs.Dozzle.Subdomain = "logs"
		obs.Dozzle.DataVolume = "dozzle_data"

		protectDozzle, err := c.prompter.Confirm("Protect Dozzle with password?", true)
		if err != nil {
			return obs, err
		}

		if protectDozzle {
			obs.Dozzle.BasicAuth, err = c.configureBasicAuth("Dozzle")
			if err != nil {
				return obs, err
			}
		}
	}

	// Beszel Configuration
	enableBeszel, err := c.prompter.Confirm("Enable Beszel (monitoring)?", true)
	if err != nil {
		return obs, err
	}

	if enableBeszel {
		obs.Beszel.Enabled = true
		obs.Beszel.Subdomain = "monitor"
		obs.Beszel.DataVolume = "beszel_data"
		obs.Beszel.SocketVolume = "beszel_socket"

		protectBeszel, err := c.prompter.Confirm("Protect Beszel with password?", true)
		if err != nil {
			return obs, err
		}

		if protectBeszel {
			obs.Beszel.BasicAuth, err = c.configureBasicAuth("Beszel")
			if err != nil {
				return obs, err
			}
		}
	}

	// Docker Socket Configuration
	customSocket, err := c.prompter.Confirm("Use custom Docker socket path?", false)
	if err != nil {
		return obs, err
	}

	if customSocket {
		socketPath, err := c.prompter.Text("Enter Docker socket path", "/var/run/docker.sock")
		if err != nil {
			return obs, err
		}
		obs.DockerSocket = socketPath
	}

	return obs, nil
}

func (c *SetupCommand) configureBasicAuth(serviceName string) (*config.BasicAuth, error) {
	c.output.Infof("🔐 Configuring authentication for %s", serviceName)

	authType, err := c.prompter.Select("Choose authentication type", []string{
		"Single user (username + password)",
		"Multiple users",
		"Generate random password",
	}, 0)
	if err != nil {
		return nil, err
	}

	auth := &config.BasicAuth{Enabled: true}

	switch {
	case strings.Contains(authType, "Single user"):
		username, err := c.prompter.Text("Enter username", "admin")
		if err != nil {
			return nil, err
		}

		password, err := c.prompter.PasswordWithValidation("Enter password", prompt.ValidatePassword)
		if err != nil {
			return nil, err
		}

		hashedPassword, err := c.hashPassword(password)
		if err != nil {
			return nil, err
		}

		auth.Username = username
		auth.Password = hashedPassword

	case strings.Contains(authType, "Multiple users"):
		auth.Users = make(map[string]string)

		for {
			username, err := c.prompter.Text("Enter username (or press Enter to finish)")
			if err != nil {
				return nil, err
			}
			if username == "" {
				break
			}

			password, err := c.prompter.PasswordWithValidation(fmt.Sprintf("Enter password for %s", username), prompt.ValidatePassword)
			if err != nil {
				return nil, err
			}

			hashedPassword, err := c.hashPassword(password)
			if err != nil {
				return nil, err
			}

			auth.Users[username] = hashedPassword
		}

	case strings.Contains(authType, "Generate random"):
		password := c.generateRandomPassword()
		hashedPassword, err := c.hashPassword(password)
		if err != nil {
			return nil, err
		}

		auth.Username = "admin"
		auth.Password = hashedPassword

		c.output.Infof("🔑 Generated credentials for %s:", serviceName)
		c.output.Infof("   Username: admin")
		c.output.Infof("   Password: %s", password)
		c.output.Info("   ⚠️  Save these credentials - they won't be shown again!")
	}

	return auth, nil
}

func (c *SetupCommand) buildStack(domain, email, project, env string, observability config.Observability) *config.Stack {
	stack := &config.Stack{
		Version: 1,
		Project: project,
		Domain:  domain,
		TLS: config.TLS{
			Mode:     "acme",
			Email:    email,
			Resolver: "le",
		},
		Observability: observability,
		Networks: map[string]config.Network{
			"private": {Internal: true},
			"public":  {Internal: false},
		},
		Volumes: []config.Volume{
			{Name: "traefik_acme"},
		},
		Services: []config.Service{},
	}

	// Add volumes for enabled services
	if observability.Dozzle.Enabled {
		stack.Volumes = append(stack.Volumes, config.Volume{Name: observability.Dozzle.DataVolume})
	}
	if observability.Beszel.Enabled {
		stack.Volumes = append(stack.Volumes,
			config.Volume{Name: observability.Beszel.DataVolume},
			config.Volume{Name: observability.Beszel.SocketVolume},
		)
	}

	// Adjust for local environment
	if env == "local" {
		stack.TLS.Mode = "disabled"
		stack.Domain = "localhost"
	}

	return stack
}

func (c *SetupCommand) showSummary(stack *config.Stack, filename string) {
	c.output.Info("")
	c.output.Info("✅ Configuration created successfully!")
	c.output.Info("")
	c.output.Infof("📄 File: %s", filename)
	c.output.Infof("🌐 Domain: %s", stack.Domain)
	c.output.Infof("📧 Email: %s", stack.TLS.Email)
	c.output.Infof("🔒 SSL: %s", stack.TLS.Mode)

	if stack.Observability.Dozzle.Enabled {
		c.output.Infof("📋 Logs: https://%s.%s", stack.Observability.Dozzle.Subdomain, stack.Domain)
	}
	if stack.Observability.Beszel.Enabled {
		c.output.Infof("📊 Monitoring: https://%s.%s", stack.Observability.Beszel.Subdomain, stack.Domain)
	}

	c.output.Info("")
	c.output.Info("🚀 Next steps:")
	c.output.Infof("   harborctl up -f %s", filename)
	c.output.Info("   harborctl status")
}

func (c *SetupCommand) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (c *SetupCommand) generateRandomPassword() string {
	bytes := make([]byte, 12)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:12]
}
