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

// gitOutputLines runs a git command and returns non-empty output lines.
func gitOutputLines(args ...string) ([]string, error) {
	output, err := exec.Command("git", args...).Output()
	if err != nil {
		return nil, err
	}
	var lines []string
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

// GetUnstagedFiles returns a list of all unstaged files (modified tracked files and untracked files)
func (c *Client) GetUnstagedFiles() ([]string, error) {
	modified, err := gitOutputLines("diff", "--name-only")
	if err != nil {
		return nil, fmt.Errorf("failed to get modified files: %w", err)
	}

	untracked, err := gitOutputLines("ls-files", "--others", "--exclude-standard")
	if err != nil {
		return nil, fmt.Errorf("failed to get untracked files: %w", err)
	}

	return append(modified, untracked...), nil
}

// GetGitRepoRoot returns the root directory of the current git repository
func (c *Client) GetGitRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}
