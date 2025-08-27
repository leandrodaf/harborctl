package compose

// ComposeFile representa um arquivo docker-compose
type ComposeFile struct {
	Version  string                    `yaml:"version"`
	Services map[string]map[string]any `yaml:"services"`
	Networks map[string]map[string]any `yaml:"networks"`
	Volumes  map[string]map[string]any `yaml:"volumes"`
	Secrets  map[string]map[string]any `yaml:"secrets,omitempty"`
}
