package compose

import (
	"fmt"

	"github.com/leandrodaf/harborctl/internal/config"
)

// HealthCheckerImpl implementa HealthChecker
type HealthCheckerImpl struct{}

// NewHealthChecker cria um novo HealthChecker
func NewHealthChecker() HealthChecker {
	return &HealthCheckerImpl{}
}

// Build constrói a configuração de health check
func (h *HealthCheckerImpl) Build(healthConfig config.HealthCheck, port int) map[string]interface{} {
	if !healthConfig.Enabled {
		return nil
	}

	check := make(map[string]interface{})

	// Path padrão se não especificado
	path := healthConfig.Path
	if path == "" {
		path = "/health"
	}

	// Comando de health check
	if port > 0 {
		check["test"] = []string{"CMD-SHELL", fmt.Sprintf("curl -f http://localhost:%d%s || exit 1", port, path)}
	} else {
		check["test"] = []string{"CMD-SHELL", fmt.Sprintf("curl -f http://localhost%s || exit 1", path)}
	}

	// Intervalo
	if healthConfig.Interval != "" {
		check["interval"] = healthConfig.Interval
	} else {
		check["interval"] = "30s"
	}

	// Timeout
	if healthConfig.Timeout != "" {
		check["timeout"] = healthConfig.Timeout
	} else {
		check["timeout"] = "10s"
	}

	// Retries
	if healthConfig.Retries > 0 {
		check["retries"] = healthConfig.Retries
	} else {
		check["retries"] = 3
	}

	// Start period
	check["start_period"] = "60s"

	return check
}

// DeployStrategyImpl implementa DeployStrategy
type DeployStrategyImpl struct{}

// NewDeployStrategy cria um novo DeployStrategy
func NewDeployStrategy() DeployStrategy {
	return &DeployStrategyImpl{}
}

// Build constrói a configuração de deploy
func (d *DeployStrategyImpl) Build(deployConfig config.DeployConfig, replicas int) map[string]interface{} {
	deploy := make(map[string]interface{})

	// Configuração de update
	updateConfig := make(map[string]interface{})

	switch deployConfig.Strategy {
	case "recreate":
		updateConfig["order"] = "stop-first"
		updateConfig["parallelism"] = 0
	case "rolling", "":
		// Rolling update é o padrão
		updateConfig["order"] = "start-first"
		updateConfig["parallelism"] = 1
		updateConfig["delay"] = "10s"
		updateConfig["failure_action"] = "rollback"
		updateConfig["monitor"] = "60s"
		updateConfig["max_failure_ratio"] = 0.3
	}

	deploy["update_config"] = updateConfig

	// Configuração de restart
	restartPolicy := make(map[string]interface{})
	restartPolicy["condition"] = "on-failure"
	restartPolicy["delay"] = "5s"
	restartPolicy["max_attempts"] = 3
	deploy["restart_policy"] = restartPolicy

	// Réplicas se especificado
	if replicas > 1 {
		deploy["replicas"] = replicas
	}

	return deploy
}
