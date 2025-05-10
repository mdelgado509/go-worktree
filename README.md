# Go Worktree Manager

A simple command-line tool to streamline git worktree workflows for task-based development.

## Features

- Create new worktrees for specific tickets or tasks
- Delete worktrees when they're no longer needed
- List all active worktrees
- Easily navigate between different worktrees
- Consistent terminal output with color-coding
- Helpful error messages and instructions

## Installation

```bash
# Clone the repository
git clone https://github.com/mdelgado509/go-worktree.git
cd go-worktree

# Build and install
make install

# Verify installation
go-worktree --version
```

## Usage Examples

### Creating Worktrees

Create a worktree for a task/ticket using the main branch as base:

```bash
go-worktree create TICKET-123
```

Create a worktree using a different base branch:

```bash
go-worktree create TICKET-123 develop
```

You can also use the `add` or `new` aliases:

```bash
go-worktree add TICKET-123
go-worktree new TICKET-123 feature-branch
```

### Listing Worktrees

List all worktrees you have created:

```bash
go-worktree list
```

Or use the shorter alias:

```bash
go-worktree ls
```

### Navigating to Worktrees

To navigate to a worktree, use:

```bash
eval $(go-worktree cd TICKET-123)
```

This works by having the `cd` command output a shell-executable command that the `eval` then executes.

### Deleting Worktrees

Delete a worktree but keep the branch:

```bash
go-worktree delete TICKET-123
```

Delete both the worktree and the branch:

```bash
go-worktree delete TICKET-123 -d
```

You can also use aliases:

```bash
go-worktree rm TICKET-123
go-worktree remove TICKET-123 -d
```

## Organization

This tool organizes worktrees by placing them in a directory structure under `~/worktrees`:

```
~/worktrees/
├── repo1/
│   ├── TICKET-123/
│   └── TICKET-456/
└── repo2/
    └── FEATURE-789/
```

Each repository gets its own directory, and within that, each ticket/task gets its own directory containing the worktree.

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
├── integration_test.sh   # Integration test script
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── README.md             # Project documentation
└── Makefile              # Build and test automation
```

## Development

To contribute to this project:

1. Clone the repository
2. Make your changes
3. Format and check your code: `make fmt`
4. Run the tests: `make test`
5. Run the integration tests: `./integration_test.sh`
   - The integration test automatically creates a temporary test repository
   - Tests the full create-list-delete workflow
   - Cleans up all test resources on completion
6. Submit a pull request

## License

This project is released under the [MIT License](LICENSE).
