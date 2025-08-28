package ssh

import (
	"fmt"
	"strings"
)

type commandBuilder struct{}

// NewCommandBuilder creates a new command builder
func NewCommandBuilder() CommandBuilder {
	return &commandBuilder{}
}

func (c *commandBuilder) BuildLogsCommand(composePath, service string, follow bool, tail int) string {
	cmd := fmt.Sprintf("docker compose -f %s logs", composePath)

	if follow {
		cmd += " -f"
	}

	if tail > 0 {
		cmd += fmt.Sprintf(" --tail %d", tail)
	}

	if service != "" {
		cmd += fmt.Sprintf(" %s", service)
	}

	return cmd
}

func (c *commandBuilder) BuildControlCommand(action, composePath, service string, verbose bool) string {
	var commands []string

	switch action {
	case "status":
		cmd := fmt.Sprintf("docker compose -f %s ps", composePath)
		if service != "" {
			cmd += fmt.Sprintf(" %s", service)
		}
		commands = append(commands, cmd)

		if verbose {
			commands = append(commands, "echo '--- RESOURCES ---'")
			commands = append(commands, "docker stats --no-stream --format 'table {{.Name}}\\t{{.CPUPerc}}\\t{{.MemUsage}}\\t{{.MemPerc}}'")
		}

	case "restart":
		if service != "" {
			commands = append(commands, fmt.Sprintf("echo '🔄 Restarting service: %s'", service))
			commands = append(commands, fmt.Sprintf("docker compose -f %s restart %s", composePath, service))
		} else {
			commands = append(commands, "echo '🔄 Restarting all services'")
			commands = append(commands, fmt.Sprintf("docker compose -f %s restart", composePath))
		}
		commands = append(commands, fmt.Sprintf("docker compose -f %s ps", composePath))

	case "stop":
		if service != "" {
			commands = append(commands, fmt.Sprintf("echo '⏹️  Stopping service: %s'", service))
			commands = append(commands, fmt.Sprintf("docker compose -f %s stop %s", composePath, service))
		} else {
			commands = append(commands, "echo '⏹️  Stopping all services'")
			commands = append(commands, fmt.Sprintf("docker compose -f %s stop", composePath))
		}
		commands = append(commands, fmt.Sprintf("docker compose -f %s ps", composePath))

	case "start":
		if service != "" {
			commands = append(commands, fmt.Sprintf("echo '▶️  Starting service: %s'", service))
			commands = append(commands, fmt.Sprintf("docker compose -f %s start %s", composePath, service))
		} else {
			commands = append(commands, "echo '▶️  Starting all services'")
			commands = append(commands, fmt.Sprintf("docker compose -f %s start", composePath))
		}
		commands = append(commands, fmt.Sprintf("docker compose -f %s ps", composePath))

	case "details":
		if service != "" {
			commands = append(commands, fmt.Sprintf("echo '📊 Service details: %s'", service))
			commands = append(commands, fmt.Sprintf("docker compose -f %s ps %s", composePath, service))
			commands = append(commands, "echo '--- RECENT LOGS ---'")
			commands = append(commands, fmt.Sprintf("docker compose -f %s logs --tail 20 %s", composePath, service))
			commands = append(commands, "echo '--- RESOURCES ---'")
			commands = append(commands, fmt.Sprintf("docker stats %s --no-stream --format 'table {{.Name}}\\t{{.CPUPerc}}\\t{{.MemUsage}}\\t{{.MemPerc}}\\t{{.NetIO}}\\t{{.BlockIO}}'", service))
		} else {
			commands = append(commands, "echo '📊 All services details'")
			commands = append(commands, fmt.Sprintf("docker compose -f %s ps", composePath))
			commands = append(commands, "echo '--- GENERAL RESOURCES ---'")
			commands = append(commands, "docker stats --no-stream --format 'table {{.Name}}\\t{{.CPUPerc}}\\t{{.MemUsage}}\\t{{.MemPerc}}'")
		}

		case "health":
		commands = append(commands, "echo '🏥 Checking services health'")
		commands = append(commands, fmt.Sprintf("docker compose -f %s ps", composePath))
		commands = append(commands, "echo '--- CONTAINER HEALTH ---'")
		commands = append(commands, "docker ps --format 'table {{.Names}}\\t{{.Status}}\\t{{.Ports}}'")
		commands = append(commands, "echo '--- SYSTEM RESOURCES ---'")
		commands = append(commands, "docker stats --no-stream --format 'table {{.Name}}\\t{{.CPUPerc}}\\t{{.MemUsage}}\\t{{.MemPerc}}'")
		
		if service != "" {
			commands = append(commands, "echo '--- SERVICE LOGS ---'")
			commands = append(commands, fmt.Sprintf("docker compose -f %s logs --tail 10 %s", composePath, service))
		}
	}

	return strings.Join(commands, " && ")
}
