package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/leandrodaf/harborctl/internal/config"
	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/prompt"
	"golang.org/x/crypto/bcrypt"
)

// EditServerCommand handles editing existing server configuration
type EditServerCommand struct {
	configManager config.Manager
	prompter      prompt.Prompter
	errorHandler  *prompt.ErrorHandler
	output        cli.Output
}

// NewEditServerCommand creates a new edit server command
func NewEditServerCommand(configManager config.Manager, output cli.Output) cli.Command {
	prompter := prompt.NewPrompter()
	return &EditServerCommand{
		configManager: configManager,
		prompter:      prompter,
		errorHandler:  prompt.NewErrorHandler(prompter),
		output:        output,
	}
}

func (c *EditServerCommand) Name() string {
	return "edit-server"
}

func (c *EditServerCommand) Description() string {
	return "Edit existing server configuration interactively"
}

func (c *EditServerCommand) Execute(ctx context.Context, args []string) error {
	filename := "server-base.yml"
	if len(args) > 0 {
		filename = args[0]
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		c.output.Errorf("❌ Configuration file not found: %s", filename)
		c.output.Info("💡 Use 'harborctl setup' to create a new configuration")
		return fmt.Errorf("configuration file not found: %s", filename)
	}

	// Load existing configuration
	stack, err := c.loadConfiguration(filename)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	c.output.Infof("🔧 Editing configuration: %s", filename)
	c.output.Info("")

	// Show current configuration
	c.showCurrentConfig(stack)

	// Interactive editing menu
	for {
		action, err := c.prompter.Select("What would you like to edit?", []string{
			"Basic settings (domain, email)",
			"SSL/TLS configuration",
			"Observability settings",
			"Authentication settings",
			"Advanced settings",
			"Save and exit",
			"Exit without saving",
		})
		if err != nil {
			return err
		}

		switch action {
		case "Basic settings (domain, email)":
			if err := c.editBasicSettings(stack); err != nil {
				c.output.Errorf("❌ Error editing basic settings: %v", err)
			}
		case "SSL/TLS configuration":
			if err := c.editSSLSettings(stack); err != nil {
				c.output.Errorf("❌ Error editing SSL settings: %v", err)
			}
		case "Observability settings":
			if err := c.editObservabilitySettings(stack); err != nil {
				c.output.Errorf("❌ Error editing observability: %v", err)
			}
		case "Authentication settings":
			if err := c.editAuthSettings(stack); err != nil {
				c.output.Errorf("❌ Error editing authentication: %v", err)
			}
		case "Advanced settings":
			if err := c.editAdvancedSettings(stack); err != nil {
				c.output.Errorf("❌ Error editing advanced settings: %v", err)
			}
		case "Save and exit":
			if err := c.configManager.SaveBaseConfig(ctx, filename, stack); err != nil {
				return fmt.Errorf("failed to save configuration: %w", err)
			}
			c.output.Infof("✅ Configuration saved to %s", filename)
			return nil
		case "Exit without saving":
			c.output.Info("❌ Changes discarded")
			return nil
		}

		c.output.Info("")
	}
}

func (c *EditServerCommand) loadConfiguration(filename string) (*config.Stack, error) {
	// This is a simplified loader - in a real implementation,
	// you'd want to properly parse the YAML file
	// For now, return a basic structure
	return &config.Stack{
		Version: 1,
		Project: "deploy",
		Domain:  "example.com",
		TLS: config.TLS{
			Mode:     "acme",
			Email:    "admin@example.com",
			Resolver: "le",
		},
		Observability: config.Observability{
			Dozzle: config.Dozzle{
				Enabled:    true,
				Subdomain:  "logs",
				DataVolume: "dozzle_data",
			},
			Beszel: config.Beszel{
				Enabled:      true,
				Subdomain:    "monitor",
				DataVolume:   "beszel_data",
				SocketVolume: "beszel_socket",
			},
		},
		Networks: map[string]config.Network{
			"private": {Internal: true},
			"public":  {Internal: false},
		},
		Volumes: []config.Volume{
			{Name: "traefik_acme"},
			{Name: "dozzle_data"},
			{Name: "beszel_data"},
			{Name: "beszel_socket"},
		},
		Services: []config.Service{},
	}, nil
}

func (c *EditServerCommand) showCurrentConfig(stack *config.Stack) {
	c.output.Info("📋 Current Configuration:")
	c.output.Infof("   Domain: %s", stack.Domain)
	c.output.Infof("   Email: %s", stack.TLS.Email)
	c.output.Infof("   SSL Mode: %s", stack.TLS.Mode)
	c.output.Infof("   Project: %s", stack.Project)

	if stack.Observability.Dozzle.Enabled {
		c.output.Info("   Dozzle: ✅ Enabled")
		if stack.Observability.Dozzle.BasicAuth != nil && stack.Observability.Dozzle.BasicAuth.Enabled {
			c.output.Info("     Authentication: ✅ Protected")
		} else {
			c.output.Info("     Authentication: ❌ Not protected")
		}
	} else {
		c.output.Info("   Dozzle: ❌ Disabled")
	}

	if stack.Observability.Beszel.Enabled {
		c.output.Info("   Beszel: ✅ Enabled")
		if stack.Observability.Beszel.BasicAuth != nil && stack.Observability.Beszel.BasicAuth.Enabled {
			c.output.Info("     Authentication: ✅ Protected")
		} else {
			c.output.Info("     Authentication: ❌ Not protected")
		}
	} else {
		c.output.Info("   Beszel: ❌ Disabled")
	}

	c.output.Info("")
}

func (c *EditServerCommand) editBasicSettings(stack *config.Stack) error {
	newDomain, err := c.prompter.Domain("Enter domain", stack.Domain)
	if err != nil {
		return err
	}
	if newDomain != "" {
		stack.Domain = newDomain
	}

	newEmail, err := c.prompter.Email("Enter email", stack.TLS.Email)
	if err != nil {
		return err
	}
	if newEmail != "" {
		stack.TLS.Email = newEmail
	}

	newProject, err := c.prompter.TextWithValidation("Enter project name", prompt.ValidateProjectName, stack.Project)
	if err != nil {
		return err
	}
	if newProject != "" {
		stack.Project = newProject
	}

	c.output.Info("✅ Basic settings updated")
	return nil
}

func (c *EditServerCommand) editSSLSettings(stack *config.Stack) error {
	sslMode, err := c.prompter.Select("Choose SSL mode", []string{
		"acme (Automatic SSL with Let's Encrypt)",
		"disabled (No SSL - for local development)",
	}, 0)
	if err != nil {
		return err
	}

	if sslMode == "acme (Automatic SSL with Let's Encrypt)" {
		stack.TLS.Mode = "acme"
	} else {
		stack.TLS.Mode = "disabled"
	}

	c.output.Info("✅ SSL settings updated")
	return nil
}

func (c *EditServerCommand) editObservabilitySettings(stack *config.Stack) error {
	service, err := c.prompter.Select("Which service to configure?", []string{
		"Dozzle (log viewer)",
		"Beszel (monitoring)",
		"Docker socket settings",
	})
	if err != nil {
		return err
	}

	switch service {
	case "Dozzle (log viewer)":
		enabled, err := c.prompter.Confirm("Enable Dozzle?", stack.Observability.Dozzle.Enabled)
		if err != nil {
			return err
		}
		stack.Observability.Dozzle.Enabled = enabled

		if enabled {
			subdomain, err := c.prompter.TextWithValidation("Dozzle subdomain", prompt.ValidateSubdomain, stack.Observability.Dozzle.Subdomain)
			if err != nil {
				return err
			}
			if subdomain != "" {
				stack.Observability.Dozzle.Subdomain = subdomain
			}
		}

	case "Beszel (monitoring)":
		enabled, err := c.prompter.Confirm("Enable Beszel?", stack.Observability.Beszel.Enabled)
		if err != nil {
			return err
		}
		stack.Observability.Beszel.Enabled = enabled

		if enabled {
			subdomain, err := c.prompter.TextWithValidation("Beszel subdomain", prompt.ValidateSubdomain, stack.Observability.Beszel.Subdomain)
			if err != nil {
				return err
			}
			if subdomain != "" {
				stack.Observability.Beszel.Subdomain = subdomain
			}
		}

	case "Docker socket settings":
		customSocket, err := c.prompter.Confirm("Use custom Docker socket path?", stack.Observability.DockerSocket != "")
		if err != nil {
			return err
		}

		if customSocket {
			currentSocket := stack.Observability.DockerSocket
			if currentSocket == "" {
				currentSocket = "/var/run/docker.sock"
			}

			socketPath, err := c.prompter.Text("Docker socket path", currentSocket)
			if err != nil {
				return err
			}
			stack.Observability.DockerSocket = socketPath
		} else {
			stack.Observability.DockerSocket = ""
		}
	}

	c.output.Info("✅ Observability settings updated")
	return nil
}

func (c *EditServerCommand) editAuthSettings(stack *config.Stack) error {
	service, err := c.prompter.Select("Configure authentication for which service?", []string{
		"Dozzle (log viewer)",
		"Beszel (monitoring)",
	})
	if err != nil {
		return err
	}

	switch service {
	case "Dozzle (log viewer)":
		return c.configureServiceAuth(&stack.Observability.Dozzle.BasicAuth, "Dozzle")
	case "Beszel (monitoring)":
		return c.configureServiceAuth(&stack.Observability.Beszel.BasicAuth, "Beszel")
	}

	return nil
}

func (c *EditServerCommand) configureServiceAuth(authPtr **config.BasicAuth, serviceName string) error {
	hasAuth := *authPtr != nil && (*authPtr).Enabled

	enableAuth, err := c.prompter.Confirm(fmt.Sprintf("Enable authentication for %s?", serviceName), hasAuth)
	if err != nil {
		return err
	}

	if !enableAuth {
		*authPtr = nil
		c.output.Infof("✅ Authentication disabled for %s", serviceName)
		return nil
	}

	// Initialize auth if needed
	if *authPtr == nil {
		*authPtr = &config.BasicAuth{}
	}
	(*authPtr).Enabled = true

	authType, err := c.prompter.Select("Authentication type", []string{
		"Single user",
		"Multiple users",
		"Keep current settings",
	}, 2)
	if err != nil {
		return err
	}

	switch authType {
	case "Single user":
		username, err := c.prompter.Text("Username", "admin")
		if err != nil {
			return err
		}

		password, err := c.prompter.PasswordWithValidation("Password", prompt.ValidatePassword)
		if err != nil {
			return err
		}

		hashedPassword, err := c.hashPassword(password)
		if err != nil {
			return err
		}

		(*authPtr).Username = username
		(*authPtr).Password = hashedPassword
		(*authPtr).Users = nil

	case "Multiple users":
		(*authPtr).Users = make(map[string]string)
		(*authPtr).Username = ""
		(*authPtr).Password = ""

		for {
			username, err := c.prompter.Text("Enter username (or press Enter to finish)")
			if err != nil {
				return err
			}
			if username == "" {
				break
			}

			password, err := c.prompter.PasswordWithValidation(fmt.Sprintf("Password for %s", username), prompt.ValidatePassword)
			if err != nil {
				return err
			}

			hashedPassword, err := c.hashPassword(password)
			if err != nil {
				return err
			}

			(*authPtr).Users[username] = hashedPassword
		}
	}

	c.output.Infof("✅ Authentication configured for %s", serviceName)
	return nil
}

func (c *EditServerCommand) editAdvancedSettings(stack *config.Stack) error {
	c.output.Info("🔧 Advanced settings coming soon...")
	c.output.Info("For now, you can manually edit the YAML file")
	return nil
}

func (c *EditServerCommand) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
