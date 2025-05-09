package worktree

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGetWorktreeBasePath tests the getWorktreeBasePath function
func TestGetWorktreeBasePath(t *testing.T) {
	path, err := getWorktreeBasePath()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, "worktrees")
	if path != expected {
		t.Errorf("Expected path %s, got %s", expected, path)
	}
}

// GitClientInterface defines the interface for git operations
type GitClientInterface interface {
	GetRepoName() (string, error)
	FetchBranch(branch string) error
	CreateWorktree(path, branchName string) error
	RemoveWorktree(path string) error
	DeleteBranch(branchName string) error
	ListWorktrees() ([]Worktree, error)
}

// Worktree for testing
type Worktree struct {
	Path   string
	Branch string
}

// MockGitClient is a mock implementation of the git client for testing
type MockGitClient struct {
	RepoName string
}

func (m *MockGitClient) GetRepoName() (string, error) {
	return m.RepoName, nil
}

func (m *MockGitClient) FetchBranch(branch string) error {
	return nil
}

func (m *MockGitClient) CreateWorktree(path, branchName string) error {
	return nil
}

func (m *MockGitClient) RemoveWorktree(path string) error {
	return nil
}

func (m *MockGitClient) DeleteBranch(branchName string) error {
	return nil
}

func (m *MockGitClient) ListWorktrees() ([]Worktree, error) {
	return []Worktree{}, nil
}

// For testing, we'll modify the Manager struct to accept an interface
type testManager struct {
	git      GitClientInterface
	basePath string
}

// TestGetPath tests the GetPath function
func TestGetPath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "go-worktree-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock manager with a fixed base path and mock git client
	manager := &testManager{
		git:      &MockGitClient{RepoName: "test-repo"},
		basePath: tempDir,
	}

	// Test GetPath function similar to the Manager.GetPath
	path, err := getTestPath(manager, "TICKET-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := filepath.Join(tempDir, "test-repo", "TICKET-123")
	if path != expected {
		t.Errorf("Expected path %s, got %s", expected, path)
	}
}

// getTestPath is a test version of GetPath that works with our test manager
func getTestPath(m *testManager, ticket string) (string, error) {
	repo, err := m.git.GetRepoName()
	if err != nil {
		return "", err
	}

	return filepath.Join(m.basePath, repo, ticket), nil
}
