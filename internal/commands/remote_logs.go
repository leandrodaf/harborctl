package commands

import (
	"context"
	"flag"

	"github.com/leandrodaf/harborctl/pkg/cli"
	"github.com/leandrodaf/harborctl/pkg/ssh"
	"github.com/leandrodaf/harborctl/pkg/validation"
)

// RemoteLogsCommand handles remote log viewing
type RemoteLogsCommand struct {
	sshExecutor    ssh.Executor
	commandBuilder ssh.CommandBuilder
	validator      validation.Validator
	output         cli.Output
}

// NewRemoteLogsCommand creates a new remote logs command
func NewRemoteLogsCommand(output cli.Output) cli.Command {
	return &RemoteLogsCommand{
		sshExecutor:    ssh.NewExecutor(),
		commandBuilder: ssh.NewCommandBuilder(),
		validator:      validation.NewValidator(),
		output:         output,
	}
}

func (c *RemoteLogsCommand) Name() string {
	return "remote-logs"
}

func (c *RemoteLogsCommand) Description() string {
	return "Connect to remote server to view service logs"
}

func (c *RemoteLogsCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("remote-logs", flag.ExitOnError)

	var host, user, keyFile, service, composePath string
	var follow bool
	var tail int
	var port int

	fs.StringVar(&host, "host", "", "remote server address")
	fs.StringVar(&user, "user", "root", "SSH user")
	fs.StringVar(&keyFile, "key", "", "SSH private key file")
	fs.IntVar(&port, "port", 22, "SSH port")
	fs.StringVar(&service, "service", "", "specific service name")
	fs.StringVar(&composePath, "compose", ".deploy/compose.generated.yml", "compose file path")
	fs.BoolVar(&follow, "follow", false, "follow logs in real time")
	fs.IntVar(&tail, "tail", 100, "number of lines to show")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := c.validator.ValidateHost(host); err != nil {
		c.output.Error("❌ Host is required. Use --host")
		return err
	}

	config := ssh.Config{
		Host:    host,
		User:    user,
		KeyFile: keyFile,
		Port:    port,
	}

	command := c.commandBuilder.BuildLogsCommand(composePath, service, follow, tail)

	if service != "" {
		c.output.Infof("📋 Connecting to %s to view logs for service: %s", host, service)
	} else {
		c.output.Infof("📋 Connecting to %s to view logs for all services", host)
	}

	if follow {
		c.output.Info("🔄 Follow mode enabled - press Ctrl+C to exit")
	}

	c.output.Info("🔗 Connecting via SSH...")
	return c.sshExecutor.Execute(ctx, config, command)
}
