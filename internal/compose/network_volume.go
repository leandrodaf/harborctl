package compose

import (
	"context"

	"github.com/leandrodaf/harborctl/internal/config"
)

// networkBuilder implementa NetworkBuilder
type networkBuilder struct{}

func NewNetworkBuilder() NetworkBuilder {
	return &networkBuilder{}
}

func (b *networkBuilder) Build(ctx context.Context, networks map[string]config.Network) map[string]map[string]any {
	result := make(map[string]map[string]any)
	for name, spec := range networks {
		network := map[string]any{
			"driver": "bridge",
		}

		if spec.Internal {
			// Rede interna - sem acesso à internet
			network["internal"] = true
			network["driver_opts"] = map[string]string{
				"com.docker.network.bridge.enable_ip_masquerade": "false",
				"com.docker.network.bridge.enable_icc":           "true",
				"com.docker.network.bridge.host_binding_ipv4":    "127.0.0.1",
			}
		} else {
			// Rede pública - com acesso controlado
			network["driver_opts"] = map[string]string{
				"com.docker.network.bridge.enable_ip_masquerade": "true",
				"com.docker.network.bridge.enable_icc":           "true",
			}
		}

		result[name] = network
	}
	return result
}

// volumeBuilder implementa VolumeBuilder
type volumeBuilder struct{}

func NewVolumeBuilder() VolumeBuilder {
	return &volumeBuilder{}
}

func (b *volumeBuilder) Build(ctx context.Context, volumes []config.Volume) map[string]map[string]any {
	result := make(map[string]map[string]any)
	for _, v := range volumes {
		result[v.Name] = map[string]any{}
	}
	return result
}
