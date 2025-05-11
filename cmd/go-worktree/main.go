// Package main provides the entry point for the go-worktree command
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mdelgado509/go-worktree/internal/util"
	"github.com/mdelgado509/go-worktree/internal/worktree"
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
	fmt.Println("  go-worktree create ABC-746                      Create worktree for ticket ABC-746")
	fmt.Println("  go-worktree create ABC-746 develop              Create from develop branch")
	fmt.Println("  go-worktree delete ABC-746 -d                   Delete worktree and branch")
	fmt.Println("  eval $(go-worktree cd ABC-746)                  Switch to ABC-746 worktree")
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
