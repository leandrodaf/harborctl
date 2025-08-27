package fs

import "context"

// FileSystem abstrai operações de sistema de arquivos
type FileSystem interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm int) error
	Exists(path string) bool
	MkdirAll(path string, perm int) error
}

// ConfigLoader carrega configurações
type ConfigLoader interface {
	Load(ctx context.Context, path string) ([]byte, error)
}

// TemplateRenderer renderiza templates
type TemplateRenderer interface {
	Render(ctx context.Context, template string, data interface{}) ([]byte, error)
}
