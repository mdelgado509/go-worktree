package git

import (
	"os"
	"os/exec"
	"testing"
)

// TestGetRepoName tests the GetRepoName function
func TestGetRepoName(t *testing.T) {
	// Skip if not in a git repository
	if _, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output(); err != nil {
		t.Skip("Skipping test: not in a git repository")
	}

	client := NewClient()
	name, err := client.GetRepoName()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if name == "" {
		t.Errorf("Expected non-empty repository name")
	}
}

// TestListWorktrees tests the ListWorktrees function
func TestListWorktrees(t *testing.T) {
	// Skip if not in a git repository
	if _, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output(); err != nil {
		t.Skip("Skipping test: not in a git repository")
	}

	client := NewClient()
	worktrees, err := client.ListWorktrees()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// We should have at least one worktree (the current one)
	if len(worktrees) < 1 {
		t.Errorf("Expected at least one worktree")
	}
}

// TestIntegration tests creating and removing a worktree
// This is more of an integration test and will modify your git repository
func TestIntegration(t *testing.T) {
	// Skip by default since this is destructive
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test. Set RUN_INTEGRATION_TESTS=1 to run")
	}

	// Skip if not in a git repository
	if _, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output(); err != nil {
		t.Skip("Skipping test: not in a git repository")
	}

	client := NewClient()
	testBranch := "test-worktree-branch"
	testPath := "/tmp/test-worktree"

	// Clean up any previous test remnants
	exec.Command("git", "worktree", "remove", testPath).Run()
	exec.Command("git", "branch", "-D", testBranch).Run()

	// Test creating a worktree
	err := client.CreateWorktree(testPath, testBranch)
	if err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}

	// Verify the worktree exists
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Errorf("Worktree directory not created")
	}

	// Test removing the worktree
	err = client.RemoveWorktree(testPath)
	if err != nil {
		t.Fatalf("Failed to remove worktree: %v", err)
	}

	// Verify the worktree is gone
	if _, err := os.Stat(testPath); !os.IsNotExist(err) {
		t.Errorf("Worktree directory still exists after removal")
	}

	// Clean up
	exec.Command("git", "branch", "-D", testBranch).Run()
}
