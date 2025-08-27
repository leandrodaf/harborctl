package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Client implementa operações git com segurança
type Client struct{}

// NewClient cria um novo cliente git
func NewClient() *Client {
	return &Client{}
}

// Clone clona um repositório
func (c *Client) Clone(ctx context.Context, url, path, token string) error {
	// Validar URL
	if err := c.validateGitURL(url); err != nil {
		return fmt.Errorf("URL inválida: %w", err)
	}

	// Preparar comando
	var cmd *exec.Cmd
	if token != "" {
		// Adicionar token para autenticação
		authenticatedURL := c.addTokenToURL(url, token)
		cmd = exec.CommandContext(ctx, "git", "clone", authenticatedURL, path)
	} else {
		cmd = exec.CommandContext(ctx, "git", "clone", url, path)
	}

	// Executar comando
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao clonar repositório: %s", string(output))
	}

	return nil
}

// Pull atualiza um repositório existente
func (c *Client) Pull(ctx context.Context, path, branch string) error {
	// Verificar se é um repositório git válido
	if !c.isValidGitRepo(path) {
		return fmt.Errorf("não é um repositório git válido: %s", path)
	}

	// Mudar para o diretório
	oldDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("erro ao acessar diretório: %w", err)
	}

	// Checkout da branch se especificada
	if branch != "" {
		cmd := exec.CommandContext(ctx, "git", "checkout", branch)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("erro ao fazer checkout da branch %s: %s", branch, string(output))
		}
	}

	// Pull
	cmd := exec.CommandContext(ctx, "git", "pull", "origin")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("erro ao fazer pull: %s", string(output))
	}

	return nil
}

// GetLatestCommit obtém o hash do último commit
func (c *Client) GetLatestCommit(ctx context.Context, path string) (string, error) {
	if !c.isValidGitRepo(path) {
		return "", fmt.Errorf("não é um repositório git válido: %s", path)
	}

	cmd := exec.CommandContext(ctx, "git", "-C", path, "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("erro ao obter commit: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// validateGitURL valida uma URL git
func (c *Client) validateGitURL(url string) error {
	// Verificar se é uma URL válida
	if url == "" {
		return fmt.Errorf("URL vazia")
	}

	// Verificar protocolos permitidos
	validPrefixes := []string{
		"https://github.com/",
		"https://gitlab.com/",
		"https://bitbucket.org/",
		"git@github.com:",
		"git@gitlab.com:",
	}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(url, prefix) {
			return nil
		}
	}

	return fmt.Errorf("URL não está na lista de permitidas")
}

// addTokenToURL adiciona token de autenticação à URL
func (c *Client) addTokenToURL(url, token string) string {
	if strings.HasPrefix(url, "https://github.com/") {
		return strings.Replace(url, "https://", fmt.Sprintf("https://%s@", token), 1)
	}
	if strings.HasPrefix(url, "https://gitlab.com/") {
		return strings.Replace(url, "https://", fmt.Sprintf("https://oauth2:%s@", token), 1)
	}
	return url
}

// isValidGitRepo verifica se um diretório é um repositório git válido
func (c *Client) isValidGitRepo(path string) bool {
	gitDir := filepath.Join(path, ".git")
	if stat, err := os.Stat(gitDir); err == nil {
		return stat.IsDir()
	}
	return false
}
