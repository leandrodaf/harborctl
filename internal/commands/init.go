package commands

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/leandrodaf/harborctl/internal/config"
	"github.com/leandrodaf/harborctl/internal/crypto"
	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/prompt"
	"golang.org/x/crypto/bcrypt"
)

type initCommand struct {
	configManager config.Manager
	prompter      prompt.Prompter
	errorHandler  *prompt.ErrorHandler
	output        cli.Output
}

func NewInitCommand(configManager config.Manager, output cli.Output) cli.Command {
	prompter := prompt.NewPrompter()
	return &initCommand{
		configManager: configManager,
		prompter:      prompter,
		errorHandler:  prompt.NewErrorHandler(prompter),
		output:        output,
	}
}

func (c *initCommand) Name() string {
	return "init"
}

func (c *initCommand) Description() string {
	return "Initialize project configuration (interactive or direct mode)"
}

func (c *initCommand) Execute(ctx context.Context, args []string) error {
	defer c.errorHandler.RecoverFromPanic()

	fs := flag.NewFlagSet("init", flag.ExitOnError)

	var domain, email, project, env string
	var noDozzle, noBeszel, interactive, help bool

	fs.StringVar(&domain, "domain", "", "base domain (ex: example.com)")
	fs.StringVar(&email, "email", "", "email for ACME certificates")
	fs.StringVar(&project, "project", "app", "project name")
	fs.StringVar(&env, "env", "", "environment (local|production)")
	fs.BoolVar(&interactive, "interactive", false, "use interactive mode")
	fs.BoolVar(&noDozzle, "no-dozzle", false, "don't include dozzle")
	fs.BoolVar(&noBeszel, "no-beszel", false, "don't include beszel")
	fs.BoolVar(&help, "help", false, "show help")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if help {
		c.showHelp()
		return nil
	}

	// Check if we should use interactive mode
	useInteractive := interactive || domain == ""

	if useInteractive {
		return c.errorHandler.SafeOperation(ctx, "Interactive Project Setup", func() error {
			return c.runInteractiveSetup(ctx)
		})
	}

	// Use direct flags mode
	return c.errorHandler.SafeOperation(ctx, "Direct Project Setup", func() error {
		return c.runDirectSetup(ctx, domain, email, project, env, noDozzle, noBeszel)
	})
}

func (c *initCommand) showHelp() {
	c.output.Info("🚀 HarborCtl Init Command")
	c.output.Info("")
	c.output.Info("USAGE:")
	c.output.Info("  harborctl init [flags]")
	c.output.Info("  harborctl init --interactive")
	c.output.Info("")
	c.output.Info("FLAGS:")
	c.output.Info("  --domain <domain>    Base domain (ex: example.com)")
	c.output.Info("  --email <email>      Email for ACME certificates")
	c.output.Info("  --project <name>     Project name (default: app)")
	c.output.Info("  --env <env>          Environment (local|production)")
	c.output.Info("  --no-dozzle          Don't include Dozzle (log viewer)")
	c.output.Info("  --no-beszel          Don't include Beszel (monitoring)")
	c.output.Info("  --interactive        Use interactive mode")
	c.output.Info("  --help               Show this help")
	c.output.Info("")
	c.output.Info("EXAMPLES:")
	c.output.Info("  harborctl init --interactive")
	c.output.Info("  harborctl init --domain example.com --email admin@example.com")
	c.output.Info("  harborctl init --domain localhost --env local")
	c.output.Info("")
}

func (c *initCommand) runInteractiveSetup(ctx context.Context) error {
	c.output.Info("🚀 Welcome to HarborCtl Project Setup!")
	c.output.Info("This wizard will help you create your stack configuration.")
	c.output.Info("")

	// Step 1: Project Name
	project, err := c.prompter.TextWithValidation(
		"What's your project name?",
		prompt.CombineValidators(
			prompt.ValidateRequired,
			prompt.ValidateProjectName,
		),
		"app",
	)
	if err != nil {
		return fmt.Errorf("failed to get project name: %w", err)
	}

	// Step 2: Environment
	env, err := c.prompter.Select(
		"Choose your environment",
		[]string{
			"Local Development",
			"Production",
		},
		0,
	)
	if err != nil {
		return fmt.Errorf("failed to get environment: %w", err)
	}

	// Convert to simple value
	if strings.Contains(env, "Local") {
		env = "local"
	} else {
		env = "production"
	}

	// Step 3: Domain
	var domain string
	if env == "local" {
		domain = "localhost"
		useCustomDomain, err := c.prompter.Confirm("Use custom domain for local development?", false)
		if err != nil {
			return fmt.Errorf("failed to get domain preference: %w", err)
		}
		if useCustomDomain {
			domain, err = c.prompter.Domain("Enter your local domain", "localhost")
			if err != nil {
				return fmt.Errorf("failed to get domain: %w", err)
			}
		}
	} else {
		domain, err = c.prompter.Domain("Enter your domain (e.g., example.com)")
		if err != nil {
			return fmt.Errorf("failed to get domain: %w", err)
		}
	}

	// Step 4: Email (if production)
	var email string
	if env == "production" {
		email, err = c.prompter.Email("Enter your email for SSL certificates", fmt.Sprintf("admin@%s", domain))
		if err != nil {
			return fmt.Errorf("failed to get email: %w", err)
		}
	}

	// Step 5: Observability services
	includeDozzle, err := c.prompter.Confirm("Include Dozzle (log viewer)?", true)
	if err != nil {
		return fmt.Errorf("failed to get Dozzle preference: %w", err)
	}

	var dozzleAuth bool
	var dozzleUsername, dozzlePassword string
	if includeDozzle {
		dozzleAuth, err = c.prompter.Confirm("Protect Dozzle with authentication?", true)
		if err != nil {
			return fmt.Errorf("failed to get Dozzle auth preference: %w", err)
		}

		if dozzleAuth {
			dozzleUsername, err = c.prompter.Text("Dozzle username", "admin")
			if err != nil {
				return fmt.Errorf("failed to get Dozzle username: %w", err)
			}

			dozzlePasswordRaw, err := c.prompter.PasswordWithValidation("Dozzle password", prompt.ValidatePassword)
			if err != nil {
				return fmt.Errorf("failed to get Dozzle password: %w", err)
			}

			dozzlePassword, err = c.hashPassword(dozzlePasswordRaw)
			if err != nil {
				return fmt.Errorf("failed to hash Dozzle password: %w", err)
			}
		}
	}

	includeBeszel, err := c.prompter.Confirm("Include Beszel (monitoring)?", true)
	if err != nil {
		return fmt.Errorf("failed to get Beszel preference: %w", err)
	}

	// Beszel usa sistema de autenticação próprio - não configurar basic auth
	if includeBeszel {
		c.output.Info("💡 Beszel uses its own built-in authentication system")
	}

	// Generate Beszel keys automatically if Beszel is enabled
	var beszelHubKey, beszelToken string
	if includeBeszel {
		c.output.Info("🔐 Generating Beszel authentication keys...")

		// Generate SSH key pair for Beszel authentication
		pubKey, _, err := crypto.GenerateED25519KeyPair()
		if err != nil {
			return fmt.Errorf("failed to generate Beszel SSH keys: %w", err)
		}
		beszelHubKey = pubKey

		// Generate token for Beszel authentication
		token, err := crypto.GenerateBeszelToken()
		if err != nil {
			return fmt.Errorf("failed to generate Beszel token: %w", err)
		}
		beszelToken = token

		c.output.Info("✅ Beszel authentication keys generated successfully")
	}

	// Create configuration
	options := config.CreateOptions{
		Domain:         domain,
		Email:          email,
		Project:        project,
		Environment:    env,
		NoDozzle:       !includeDozzle,
		NoBeszel:       !includeBeszel,
		DozzleAuth:     dozzleAuth,
		BeszelAuth:     false, // Sempre false - Beszel usa auth próprio
		DozzleUsername: dozzleUsername,
		DozzlePassword: dozzlePassword,
		BeszelUsername: "", // Não usado
		BeszelPassword: "", // Não usado
		BeszelHubKey:   beszelHubKey,
		BeszelToken:    beszelToken,
	}

	// Show summary
	c.showConfigSummary(project, domain, email, env, includeDozzle, includeBeszel, dozzleAuth, false)

	confirm, err := c.prompter.Confirm("Create project with these settings?", true)
	if err != nil {
		return fmt.Errorf("failed to confirm creation: %w", err)
	}

	if !confirm {
		c.output.Info("Project creation cancelled.")
		return nil
	}

	return c.createProject(ctx, options)
}

func (c *initCommand) runDirectSetup(ctx context.Context, domain, email, project, env string, noDozzle, noBeszel bool) error {
	// Set defaults
	if env == "" {
		if domain == "localhost" || domain == "" || strings.HasSuffix(domain, ".local") {
			env = "local"
		} else {
			env = "production"
		}
	}

	if domain == "" {
		if env == "local" {
			domain = "localhost"
		} else {
			return fmt.Errorf("domain is required for production environment")
		}
	}

	// Validate required fields for production
	if env == "production" && email == "" {
		return fmt.Errorf("email is required for production environment")
	}

	// Generate Beszel keys automatically if Beszel is enabled
	var beszelHubKey, beszelToken string
	if !noBeszel {
		c.output.Info("🔐 Generating Beszel authentication keys...")

		// Generate SSH key pair for Beszel authentication
		pubKey, _, err := crypto.GenerateED25519KeyPair()
		if err != nil {
			return fmt.Errorf("failed to generate Beszel SSH keys: %w", err)
		}
		beszelHubKey = pubKey

		// Generate token for Beszel authentication
		token, err := crypto.GenerateBeszelToken()
		if err != nil {
			return fmt.Errorf("failed to generate Beszel token: %w", err)
		}
		beszelToken = token

		c.output.Info("✅ Beszel authentication keys generated successfully")
	}

	options := config.CreateOptions{
		Domain:         domain,
		Email:          email,
		Project:        project,
		Environment:    env,
		NoDozzle:       noDozzle,
		NoBeszel:       noBeszel,
		BeszelAuth:     false, // Sempre false - Beszel usa auth próprio
		BeszelUsername: "",    // Não usado
		BeszelPassword: "",    // Não usado
		BeszelHubKey:   beszelHubKey,
		BeszelToken:    beszelToken,
	}

	return c.createProject(ctx, options)
}

func (c *initCommand) createProject(ctx context.Context, options config.CreateOptions) error {
	filename := "stack.yml"

	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		overwrite, err := c.prompter.Confirm(
			fmt.Sprintf("File %s already exists. Overwrite?", filename),
			false,
		)
		if err != nil {
			return fmt.Errorf("failed to confirm overwrite: %w", err)
		}

		if !overwrite {
			c.output.Info("Project creation cancelled.")
			return nil
		}
	}

	if err := c.configManager.Create(ctx, filename, options); err != nil {
		return fmt.Errorf("error creating %s: %w", filename, err)
	}

	// Create deploy directory
	if err := os.MkdirAll(".deploy", 0755); err != nil {
		return fmt.Errorf("error creating .deploy directory: %w", err)
	}

	c.output.Info("✅ Project created successfully!")
	c.output.Info(fmt.Sprintf("📄 Configuration file: %s", filename))

	// Add Beszel-specific information if enabled
	if !options.NoBeszel {
		c.output.Info("")
		c.output.Info("🔐 Beszel Authentication:")
		c.output.Info("   ✅ SSH keys and token auto-generated")
		c.output.Info("   ✅ Hub-Agent communication configured")
		c.output.Info("   ✅ Ready for immediate deployment")
	}

	c.output.Info("")
	c.output.Info("🚀 Next steps:")
	c.output.Info("   1. Edit stack.yml to add your services")
	c.output.Info("   2. Run: harborctl up")
	c.output.Info("   3. Check status: harborctl status")

	if !options.NoBeszel {
		c.output.Info("   4. Access monitoring: https://monitor." + options.Domain)
	}

	return nil
}

func (c *initCommand) showConfigSummary(project, domain, email, env string, includeDozzle, includeBeszel, dozzleAuth, beszelAuth bool) {
	c.output.Info("")
	c.output.Info("📋 Configuration Summary:")
	c.output.Info(fmt.Sprintf("   Project: %s", project))
	c.output.Info(fmt.Sprintf("   Domain: %s", domain))
	if email != "" {
		c.output.Info(fmt.Sprintf("   Email: %s", email))
	}
	c.output.Info(fmt.Sprintf("   Environment: %s", env))

	// Services
	c.output.Info("   Services:")
	if includeDozzle {
		if dozzleAuth {
			c.output.Info("     • Dozzle (log viewer): ✅ Enabled with authentication")
		} else {
			c.output.Info("     • Dozzle (log viewer): ✅ Enabled (no authentication)")
		}
	} else {
		c.output.Info("     • Dozzle (log viewer): ❌ Disabled")
	}
	if includeBeszel {
		c.output.Info("     • Beszel (monitoring): ✅ Enabled with built-in auth + keys auto-generated")
	} else {
		c.output.Info("     • Beszel (monitoring): ❌ Disabled")
	}
	c.output.Info("")
}

func (c *initCommand) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
