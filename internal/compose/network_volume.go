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
		network := map[string]any{"driver": "bridge"}
		if spec.Internal {
			network["internal"] = true
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
