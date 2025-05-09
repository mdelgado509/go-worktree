// Package git provides a client for interacting with git commands
package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Worktree represents a git worktree
type Worktree struct {
	Path   string
	Branch string
}

// Client wraps git command operations
type Client struct{}

// NewClient creates a new git client
func NewClient() *Client {
	return &Client{}
}

// GetRepoName gets the name of the current git repository
func (c *Client) GetRepoName() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not in a git repository: %w", err)
	}
	repoPath := strings.TrimSpace(string(output))
	return filepath.Base(repoPath), nil
}

// FetchBranch fetches the latest changes for a branch
func (c *Client) FetchBranch(branch string) error {
	cmd := exec.Command("git", "fetch", "origin", branch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to fetch branch: %w", err)
	}
	return nil
}

// CreateWorktree creates a new worktree with a new branch
func (c *Client) CreateWorktree(path, branchName string) error {
	cmd := exec.Command("git", "worktree", "add", path, "-b", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w", string(output), err)
	}
	return nil
}

// RemoveWorktree removes a worktree
func (c *Client) RemoveWorktree(path string) error {
	cmd := exec.Command("git", "worktree", "remove", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w", string(output), err)
	}
	return nil
}

// DeleteBranch deletes a branch
func (c *Client) DeleteBranch(branchName string) error {
	cmd := exec.Command("git", "branch", "-D", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w", string(output), err)
	}
	return nil
}

// ListWorktrees returns a list of all worktrees for the current repository
func (c *Client) ListWorktrees() ([]Worktree, error) {
	cmd := exec.Command("git", "worktree", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	// Parse output
	worktreeLines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var worktrees []Worktree

	for _, line := range worktreeLines {
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		path := parts[0]
		branch := strings.Trim(parts[2], "[]")

		worktrees = append(worktrees, Worktree{
			Path:   path,
			Branch: branch,
		})
	}

	return worktrees, nil
}
