package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Client implements secure git operations
type Client struct{}

// NewClient creates a new git client
func NewClient() *Client {
	return &Client{}
}

// Clone clones a repository
func (c *Client) Clone(ctx context.Context, url, path, token string) error {
	// Validate URL
	if err := c.validateGitURL(url); err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Prepare command
	var cmd *exec.Cmd
	if token != "" {
		// Add token for authentication
		authenticatedURL := c.addTokenToURL(url, token)
		cmd = exec.CommandContext(ctx, "git", "clone", authenticatedURL, path)
	} else {
		cmd = exec.CommandContext(ctx, "git", "clone", url, path)
	}

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error cloning repository: %s", string(output))
	}

	return nil
}

// Pull updates an existing repository
func (c *Client) Pull(ctx context.Context, path, branch string) error {
	// Check if it's a valid git repository
	if !c.isValidGitRepo(path) {
		return fmt.Errorf("not a valid git repository: %s", path)
	}

	// Change to directory
	oldDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer func() {
		if err := os.Chdir(oldDir); err != nil {
			// Log the error but don't fail the operation
			fmt.Printf("Warning: failed to restore directory: %v\n", err)
		}
	}()

	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("error accessing directory: %w", err)
	}

	// Checkout branch if specified
	if branch != "" {
		cmd := exec.CommandContext(ctx, "git", "checkout", branch)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("error checking out branch %s: %s", branch, string(output))
		}
	}

	// Pull
	cmd := exec.CommandContext(ctx, "git", "pull", "origin")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error pulling: %s", string(output))
	}

	return nil
}

// GetLatestCommit gets the latest commit hash
func (c *Client) GetLatestCommit(ctx context.Context, path string) (string, error) {
	if !c.isValidGitRepo(path) {
		return "", fmt.Errorf("not a valid git repository: %s", path)
	}

	cmd := exec.CommandContext(ctx, "git", "-C", path, "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting commit: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// validateGitURL validates a git URL
func (c *Client) validateGitURL(url string) error {
	// Check if URL is valid
	if url == "" {
		return fmt.Errorf("empty URL")
	}

	// Check allowed protocols
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

	return fmt.Errorf("URL not in allowed list")
}

// addTokenToURL adds authentication token to URL
func (c *Client) addTokenToURL(url, token string) string {
	if strings.HasPrefix(url, "https://github.com/") {
		return strings.Replace(url, "https://", fmt.Sprintf("https://%s@", token), 1)
	}
	if strings.HasPrefix(url, "https://gitlab.com/") {
		return strings.Replace(url, "https://", fmt.Sprintf("https://oauth2:%s@", token), 1)
	}
	return url
}

// isValidGitRepo checks if a directory is a valid git repository
func (c *Client) isValidGitRepo(path string) bool {
	gitDir := filepath.Join(path, ".git")
	if stat, err := os.Stat(gitDir); err == nil {
		return stat.IsDir()
	}
	return false
}
