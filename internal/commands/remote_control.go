package commands

import (
	"context"
	"flag"

	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/ssh"
	"github.com/leandrodaf/harborctl/pkg/validation"
)

// RemoteControlCommand handles remote service control
type RemoteControlCommand struct {
	sshExecutor     ssh.Executor
	commandBuilder  ssh.CommandBuilder
	validator       validation.Validator
	actionValidator validation.ActionValidator
	output          cli.Output
}

// NewRemoteControlCommand creates a new remote control command
func NewRemoteControlCommand(output cli.Output) cli.Command {
	return &RemoteControlCommand{
		sshExecutor:     ssh.NewExecutor(),
		commandBuilder:  ssh.NewCommandBuilder(),
		validator:       validation.NewValidator(),
		actionValidator: validation.NewActionValidator(),
		output:          output,
	}
}

func (c *RemoteControlCommand) Name() string {
	return "remote-control"
}

func (c *RemoteControlCommand) Description() string {
	return "Control services remotely (restart, stop, status, details)"
}

func (c *RemoteControlCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("remote-control", flag.ExitOnError)

	var host, user, keyFile, action, service, composePath string
	var port int
	var verbose bool

	fs.StringVar(&host, "host", "", "remote server address")
	fs.StringVar(&user, "user", "root", "SSH user")
	fs.StringVar(&keyFile, "key", "", "SSH private key file")
	fs.IntVar(&port, "port", 22, "SSH port")
	fs.StringVar(&action, "action", "status", "action: status, restart, stop, start, details, health")
	fs.StringVar(&service, "service", "", "specific service name")
	fs.StringVar(&composePath, "compose", ".deploy/compose.generated.yml", "compose file path")
	fs.BoolVar(&verbose, "verbose", false, "show detailed information")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := c.validator.ValidateHost(host); err != nil {
		c.output.Error("❌ Host is required. Use --host")
		return err
	}

	validActions := c.actionValidator.GetValidActions()
	if err := c.validator.ValidateAction(action, validActions); err != nil {
		c.output.Errorf("❌ %s", err.Error())
		return err
	}

	config := ssh.Config{
		Host:    host,
		User:    user,
		KeyFile: keyFile,
		Port:    port,
	}

	command := c.commandBuilder.BuildControlCommand(action, composePath, service, verbose)

	c.output.Infof("🎛️  Executing '%s' on %s", action, host)
	if service != "" {
		c.output.Infof("🎯 Service: %s", service)
	}

	c.output.Info("🔗 Connecting via SSH...")
	return c.sshExecutor.Execute(ctx, config, command)
}
