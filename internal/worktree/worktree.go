// Package worktree provides functionality for managing git worktrees
package worktree

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mdelgado509/go-worktree/internal/git"
	"github.com/mdelgado509/go-worktree/internal/util"
)

// Manager handles worktree operations
type Manager struct {
	git      *git.Client
	basePath string
}

// NewManager creates a new worktree manager
func NewManager() *Manager {
	basePath, err := getWorktreeBasePath()
	if err != nil {
		basePath = "~/worktrees" // Fallback
	}

	return &Manager{
		git:      git.NewClient(),
		basePath: basePath,
	}
}

// getWorktreeBasePath returns the base path for worktrees
func getWorktreeBasePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}
	return filepath.Join(home, "worktrees"), nil
}

// GetPath returns the path for a specific worktree
func (m *Manager) GetPath(ticket string) (string, error) {
	repo, err := m.git.GetRepoName()
	if err != nil {
		return "", err
	}

	return filepath.Join(m.basePath, repo, ticket), nil
}

// Create creates a new git worktree
func (m *Manager) Create(ticket, baseBranch string) error {
	repo, err := m.git.GetRepoName()
	if err != nil {
		return err
	}

	// Ensure base directory exists
	worktreeDir := filepath.Join(m.basePath, repo, ticket)
	err = os.MkdirAll(filepath.Dir(worktreeDir), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if directory already exists
	if _, err := os.Stat(worktreeDir); !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("directory already exists: %s", worktreeDir)
	}

	// Try to fetch latest from base branch, but don't fail if no remote exists
	fmt.Printf("Fetching latest from %s...\n", baseBranch)
	if err := m.git.FetchBranch(baseBranch); err != nil {
		fmt.Printf("Warning: couldn't fetch latest from remote (this is okay for local-only repos): %v\n", err)
	}

	// Create worktree with new branch
	fmt.Printf("Creating worktree for %s%s%s...\n", util.ColorBlue, ticket, util.ColorReset)
	if err := m.git.CreateWorktree(worktreeDir, ticket); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}

	fmt.Printf("%sSuccess!%s Worktree created at: %s\n", util.ColorGreen, util.ColorReset, worktreeDir)
	fmt.Printf("Run: %scd %s%s to start working\n", util.ColorYellow, worktreeDir, util.ColorReset)
	return nil
}

// Delete deletes a git worktree
func (m *Manager) Delete(ticket string, deleteBranch bool) error {
	worktreePath, err := m.GetPath(ticket)
	if err != nil {
		return err
	}

	// Check if worktree exists
	if _, err := os.Stat(worktreePath); errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("worktree for ticket %s not found", ticket)
	}

	// Remove worktree
	fmt.Printf("Removing worktree for %s%s%s...\n", util.ColorBlue, ticket, util.ColorReset)
	if err := m.git.RemoveWorktree(worktreePath); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	// Delete branch if requested
	if deleteBranch {
		fmt.Printf("Deleting branch %s%s%s...\n", util.ColorBlue, ticket, util.ColorReset)
		if err := m.git.DeleteBranch(ticket); err != nil {
			return fmt.Errorf("failed to delete branch: %w", err)
		}
	}

	fmt.Printf("%sDone!%s Worktree for ticket %s has been removed\n",
		util.ColorGreen, util.ColorReset, ticket)
	return nil
}

// List lists all git worktrees
func (m *Manager) List() error {
	repo, err := m.git.GetRepoName()
	if err != nil {
		return err
	}

	// Get git worktree list
	worktrees, err := m.git.ListWorktrees()
	if err != nil {
		return fmt.Errorf("failed to list worktrees: %w", err)
	}

	// Create map for easier lookup
	worktreeMap := make(map[string]string) // path -> branch
	for _, wt := range worktrees {
		worktreeMap[wt.Path] = wt.Branch
	}

	repoPath := filepath.Join(m.basePath, repo)

	// Check if directory exists
	if _, err := os.Stat(repoPath); errors.Is(err, fs.ErrNotExist) {
		fmt.Printf("No worktrees found for repository %s%s%s\n",
			util.ColorYellow, repo, util.ColorReset)
		return nil
	}

	// List directories in the repo path
	fmt.Printf("Worktrees for repository %s%s%s:\n", util.ColorYellow, repo, util.ColorReset)
	entries, err := os.ReadDir(repoPath)
	if err != nil {
		return fmt.Errorf("failed to read worktree directory: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("  No worktrees found")
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		ticket := entry.Name()
		path := filepath.Join(repoPath, ticket)
		branch, exists := worktreeMap[path]
		if !exists {
			branch = "detached"
		}

		fmt.Printf("  %s%s%s -> %s (%s%s%s)\n",
			util.ColorGreen, ticket, util.ColorReset,
			path,
			util.ColorBlue, branch, util.ColorReset)
	}

	return nil
}
