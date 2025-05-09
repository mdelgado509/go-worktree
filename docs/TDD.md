# Git Worktree Manager Technical Specification

## Overview

This document provides a detailed technical specification for a simple Git Worktree Manager CLI tool written in Go. This tool streamlines the workflow of creating, managing, and deleting git worktrees for specific tickets or tasks.

## Project Structure

```
go-worktree/
├── cmd/
│   └── go-worktree/
│       └── main.go       # Main application entry point
├── internal/
│   ├── worktree/         # Core worktree functionality
│   │   ├── worktree.go   # Worktree operations
│   │   └── worktree_test.go  # Tests for worktree operations
│   └── git/              # Git operations
│       ├── git.go        # Git command wrappers
│       └── git_test.go   # Tests for git operations
├── pkg/
│   └── util/             # Utility functions
│       ├── color.go      # Terminal color functions
│       └── color_test.go # Tests for color functions
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── README.md             # Project documentation
└── Makefile              # Build and test automation
```

## Implementation Details

### Module Definition (go.mod)

```go
module github.com/mdelgado509/go-worktree

go 1.22
```

### Main Application (cmd/go-worktree/main.go)

```go
// Package main provides the entry point for the go-worktree command
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mdelgado509/go-worktree/internal/worktree"
	"github.com/mdelgado509/go-worktree/pkg/util"
)

// Command constants define the available commands
const (
	cmdCreate = "create"
	cmdDelete = "delete"
	cmdList   = "list"
	cmdCD     = "cd"
	version   = "1.0.0"
)

// commandAliases maps alternative command names to canonical commands
var commandAliases = map[string]string{
	"new":     cmdCreate,
	"add":     cmdCreate,
	"rm":      cmdDelete,
	"cleanup": cmdDelete,
	"remove":  cmdDelete,
	"ls":      cmdList,
	"switch":  cmdCD,
}

func main() {
	// Show usage if no arguments are provided
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// Get the command and resolve aliases
	cmdArg := os.Args[1]

	// Handle special case for help and version
	if cmdArg == "help" || cmdArg == "--help" || cmdArg == "-h" {
		printUsage()
		return
	}

	if cmdArg == "version" || cmdArg == "--version" || cmdArg == "-v" {
		fmt.Printf("go-worktree version %s\n", version)
		return
	}

	// Resolve command alias
	cmd, exists := commandAliases[cmdArg]
	if !exists {
		cmd = cmdArg
	}

	// Route to appropriate handler
	switch cmd {
	case cmdCreate:
		handleCreate()
	case cmdDelete:
		handleDelete()
	case cmdList:
		handleList()
	case cmdCD:
		handleCD()
	default:
		fmt.Fprintf(os.Stderr, "%sUnknown command: %s%s\n",
			util.ColorRed, cmdArg, util.ColorReset)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Golang Git Worktree Manager - Streamlined workflow")
	fmt.Println("\nUsage:")
	fmt.Println("  go-worktree create|add TICKET-ID [BASE-BRANCH]  Create a new worktree (default: main)")
	fmt.Println("  go-worktree delete|rm TICKET-ID [-d]            Delete a worktree (-d to delete branch)")
	fmt.Println("  go-worktree list|ls                             List all your worktrees")
	fmt.Println("  go-worktree cd|switch TICKET-ID                 Print command to change to worktree")
	fmt.Println("  go-worktree help|--help                         Show this help message")
	fmt.Println("  go-worktree version|--version                   Show version information")
	fmt.Println("\nExamples:")
	fmt.Println("  go-worktree create DML-746                      Create worktree for ticket DML-746")
	fmt.Println("  go-worktree create DML-746 develop              Create from develop branch")
	fmt.Println("  go-worktree delete DML-746 -d                   Delete worktree and branch")
	fmt.Println("  eval $(go-worktree cd DML-746)                  Switch to DML-746 worktree")
}

// handleCreate handles the create command
func handleCreate() {
	createCommand := flag.NewFlagSet(cmdCreate, flag.ExitOnError)
	baseBranch := createCommand.String("base", "main", "Base branch to create from")

	// Parse remaining args
	err := createCommand.Parse(os.Args[2:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError: %v%s\n", util.ColorRed, err, util.ColorReset)
		os.Exit(1)
	}

	args := createCommand.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%sError: Ticket ID required%s\n", util.ColorRed, util.ColorReset)
		os.Exit(1)
	}

	ticket := args[0]
	// Allow overriding base branch as positional arg for convenience
	if len(args) > 1 {
		*baseBranch = args[1]
	}

	wt := worktree.NewManager()
	if err := wt.Create(ticket, *baseBranch); err != nil {
		fmt.Fprintf(os.Stderr, "%sError: %v%s\n", util.ColorRed, err, util.ColorReset)
		os.Exit(1)
	}
}

// handleDelete handles the delete command
func handleDelete() {
	deleteCommand := flag.NewFlagSet(cmdDelete, flag.ExitOnError)
	deleteBranch := deleteCommand.Bool("d", false, "Delete branch as well")

	// Parse remaining args
	err := deleteCommand.Parse(os.Args[2:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError: %v%s\n", util.ColorRed, err, util.ColorReset)
		os.Exit(1)
	}

	args := deleteCommand.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%sError: Ticket ID required%s\n", util.ColorRed, util.ColorReset)
		os.Exit(1)
	}

	ticket := args[0]
	wt := worktree.NewManager()
	if err := wt.Delete(ticket, *deleteBranch); err != nil {
		fmt.Fprintf(os.Stderr, "%sError: %v%s\n", util.ColorRed, err, util.ColorReset)
		os.Exit(1)
	}
}

// handleList handles the list command
func handleList() {
	wt := worktree.NewManager()
	if err := wt.List(); err != nil {
		fmt.Fprintf(os.Stderr, "%sError: %v%s\n", util.ColorRed, err, util.ColorReset)
		os.Exit(1)
	}
}

// handleCD handles the cd command
func handleCD() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%sError: Ticket ID required%s\n", util.ColorRed, util.ColorReset)
		os.Exit(1)
	}

	ticket := os.Args[2]
	wt := worktree.NewManager()
	path, err := wt.GetPath(ticket)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError: %v%s\n", util.ColorRed, err, util.ColorReset)
		os.Exit(1)
	}

	// Output command for shell to evaluate
	fmt.Printf("cd %s\n", path)
	fmt.Fprintf(os.Stderr, "%sNote: Run with eval $(go-worktree cd %s) to change directory%s\n",
		util.ColorYellow, ticket, util.ColorReset)
}
```

### Worktree Package (internal/worktree/worktree.go)

```go
// Package worktree provides functionality for managing git worktrees
package worktree

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mdelgado509/go-worktree/internal/git"
	"github.com/mdelgado509/go-worktree/pkg/util"
)

// Manager handles worktree operations
type Manager struct {
	git    *git.Client
	basePath string
}

// NewManager creates a new worktree manager
func NewManager() *Manager {
	basePath, err := getWorktreeBasePath()
	if err != nil {
		basePath = "~/worktrees" // Fallback
	}

	return &Manager{
		git:    git.NewClient(),
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

	// Fetch latest from base branch
	fmt.Printf("Fetching latest from %s...\n", baseBranch)
	if err := m.git.FetchBranch(baseBranch); err != nil {
		return fmt.Errorf("failed to fetch from %s: %w", baseBranch, err)
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
```

### Git Client (internal/git/git.go)

```go
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
```

### Color Utilities (pkg/util/color.go)

```go
// Package util provides utility functions
package util

// ANSI color codes for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

// Colorize returns a string with color codes
func Colorize(text, color string) string {
	return color + text + ColorReset
}

// Bold returns a string in bold
func Bold(text string) string {
	return "\033[1m" + text + "\033[0m"
}
```

## Testing Implementation

### Git Client Tests (internal/git/git_test.go)

```go
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
```

### Worktree Manager Tests (internal/worktree/worktree_test.go)

```go
package worktree

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGetPath tests the GetPath function
func TestGetPath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "go-worktree-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock manager with a fixed base path
	manager := NewManager()
	manager.basePath = tempDir

	// Mock the git client to return a fixed repo name
	oldGetRepoName := manager.git.GetRepoName
	manager.git.GetRepoName = func() (string, error) {
		return "test-repo", nil
	}
	defer func() { manager.git.GetRepoName = oldGetRepoName }()

	// Test GetPath
	path, err := manager.GetPath("TICKET-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := filepath.Join(tempDir, "test-repo", "TICKET-123")
	if path != expected {
		t.Errorf("Expected path %s, got %s", expected, path)
	}
}

// Additional tests for Create, Delete, and List would follow a similar pattern,
// creating mock git clients that return expected values and verifying behavior.
```

### Color Utilities Tests (pkg/util/color_test.go)

```go
package util

import (
	"testing"
)

// TestColorize tests the Colorize function
func TestColorize(t *testing.T) {
	testCases := []struct {
		text     string
		color    string
		expected string
	}{
		{"test", ColorRed, "\033[31mtest\033[0m"},
		{"hello", ColorBlue, "\033[34mhello\033[0m"},
		{"", ColorGreen, "\033[32m\033[0m"},
	}

	for _, tc := range testCases {
		result := Colorize(tc.text, tc.color)
		if result != tc.expected {
			t.Errorf("Expected %q, got %q", tc.expected, result)
		}
	}
}

// TestBold tests the Bold function
func TestBold(t *testing.T) {
	testCases := []struct {
		text     string
		expected string
	}{
		{"test", "\033[1mtest\033[0m"},
		{"hello", "\033[1mhello\033[0m"},
		{"", "\033[1m\033[0m"},
	}

	for _, tc := range testCases {
		result := Bold(tc.text)
		if result != tc.expected {
			t.Errorf("Expected %q, got %q", tc.expected, result)
		}
	}
}
```

## Makefile

```makefile
.PHONY: build test lint clean install

# Set build variables
BINARY_NAME=go-worktree
VERSION=1.0.0
BUILD_DIR=build
INSTALL_DIR=$(HOME)/.local/bin

# Build the binary
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/go-worktree

# Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run integration tests (more invasive)
test-integration:
	@echo "Running integration tests..."
	RUN_INTEGRATION_TESTS=1 go test -v ./...

# Run code linting
lint:
	@echo "Linting code..."
	go vet ./...
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping additional linting"; \
	fi

# Install the binary
install: build
	@echo "Installing to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@echo "Installation complete! Make sure $(INSTALL_DIR) is in your PATH."

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
```

## Development Workflow

1. Set up the project structure as outlined above
2. Implement the code in each file
3. Write unit tests for each package
4. Use the Makefile to build, test, and install the tool

## Testing Strategy

### Unit Tests

For each package, write unit tests that verify:

1. **Git Package**:

   - Test that Git commands are correctly formatted
   - Verify error handling works as expected
   - Mock external commands where appropriate

2. **Worktree Package**:

   - Test path construction
   - Verify error handling for missing repositories
   - Mock Git client responses for predictable testing

3. **Utility Package**:
   - Test color formatting functions
   - Verify string manipulation utilities

### Integration Tests

Integration tests should verify the complete workflow:

1. Create a worktree for a test ticket
2. Verify the worktree exists
3. List worktrees and verify the test ticket appears
4. Delete the worktree
5. Verify the worktree is removed

## Best Practices

1. **Error Handling**: Always wrap errors with context using `fmt.Errorf("context: %w", err)`
2. **Testing**: Write tests for all functionality
3. **Modularity**: Keep concerns separated into different packages
4. **Documentation**: Provide clear documentation for all functions and packages
5. **Command Structure**: Use the `flag` package for simple flag parsing
6. **User Experience**: Provide helpful error messages and examples
7. **Code Style**: Follow Go's official style guide and use `go fmt` and `go vet`

## Notes for Junior Engineers

1. **Go Modules**: Ensure you initialize the project with `go mod init` before starting
2. **Testing as You Go**: Write tests alongside your code, not as an afterthought
3. **Command Execution**: Be careful when executing external commands, always check for errors
4. **Error Messages**: Make error messages helpful for the end user
5. **OS Independence**: Use `os.PathSeparator` and `filepath` functions for cross-platform compatibility
6. **Git Interface**: Remember that git commands might behave differently across git versions

## Deployment

To build and install the tool:

```bash
# Clone the repository
git clone https://github.com/mdelgado509/go-worktree.git
cd go-worktree

# Build and install
make install

# Verify installation
go-worktree --version
```

This completes the technical specification for the Git Worktree Manager. The implementation provides a clean, maintainable, and testable solution that follows Go best practices.
