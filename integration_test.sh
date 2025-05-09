#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Get absolute paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMP_DIR="/tmp/go-worktree-test-$(date +%s)"
TEST_REPO_PATH="$TEMP_DIR/test-repo"
BINARY_PATH="$SCRIPT_DIR/build/go-worktree"
TEST_TICKET="TEST-123"
HOME_DIR="$HOME"
WORKTREE_BASE="$HOME_DIR/worktrees"
REPO_NAME="test-repo"
WORKTREE_PATH="$WORKTREE_BASE/$REPO_NAME/$TEST_TICKET"

# Function to clean up resources
cleanup() {
  echo -e "\n${YELLOW}Cleaning up test resources...${NC}"
  # Remove temporary directory
  if [ -d "$TEMP_DIR" ]; then
    rm -rf "$TEMP_DIR"
    echo -e "  Removed temporary directory: $TEMP_DIR"
  fi
  
  # Remove worktree if it exists
  if [ -d "$WORKTREE_PATH" ]; then
    rm -rf "$WORKTREE_PATH"
    echo -e "  Removed worktree directory: $WORKTREE_PATH"
  fi
  
  # Remove repository directory in worktrees if empty
  REPO_DIR="$WORKTREE_BASE/$REPO_NAME"
  if [ -d "$REPO_DIR" ] && [ -z "$(ls -A "$REPO_DIR")" ]; then
    rmdir "$REPO_DIR"
    echo -e "  Removed empty repository directory: $REPO_DIR"
  fi
  
  echo -e "${GREEN}Cleanup completed${NC}"
}

# Register cleanup function to run on script exit
trap cleanup EXIT

# Exit on error
set -e

# Display paths being used
echo -e "${YELLOW}Using:${NC}"
echo -e "  Script Directory: $SCRIPT_DIR"
echo -e "  Temporary Directory: $TEMP_DIR"
echo -e "  Test Repository: $TEST_REPO_PATH"
echo -e "  Binary: $BINARY_PATH"
echo -e "  Expected Worktree Path: $WORKTREE_PATH"

# Setup: Create a temporary test repository
echo -e "\n${YELLOW}=== Setting up test environment ===${NC}"
mkdir -p "$TEST_REPO_PATH"
cd "$TEST_REPO_PATH"
git init
echo "# Test Repository" > README.md
git add README.md
git config --local user.name "Test User"
git config --local user.email "test@example.com"
git commit -m "Initial commit"
echo -e "${GREEN}✓ Test repository created at: $TEST_REPO_PATH${NC}"

# Build the binary
echo -e "\n${YELLOW}Building the binary...${NC}"
cd "$SCRIPT_DIR" && make build

# Check if the binary exists
if [ ! -f "$BINARY_PATH" ]; then
  echo -e "${RED}Binary not found at $BINARY_PATH${NC}"
  exit 1
fi

# Ensure worktrees base directory exists
if [ ! -d "$WORKTREE_BASE" ]; then
  mkdir -p "$WORKTREE_BASE"
fi

# 1. Create a worktree for a test ticket
echo -e "\n${YELLOW}=== Testing worktree creation ===${NC}"

# Go to test repo
cd "$TEST_REPO_PATH"

# Create worktree using our tool
echo -e "${YELLOW}Creating worktree for $TEST_TICKET...${NC}"
"$BINARY_PATH" create "$TEST_TICKET" main

# Verify worktree exists
if [ -d "$WORKTREE_PATH" ]; then
  echo -e "${GREEN}✓ Worktree directory created successfully at $WORKTREE_PATH${NC}"
else
  echo -e "${RED}✗ Worktree directory not created at $WORKTREE_PATH${NC}"
  exit 1
fi

# 2. List worktrees
echo -e "\n${YELLOW}=== Testing worktree listing ===${NC}"
WORKTREE_LIST=$("$BINARY_PATH" list)
echo "$WORKTREE_LIST"

# Check if our test ticket is in the list
if echo "$WORKTREE_LIST" | grep -q "$TEST_TICKET"; then
  echo -e "${GREEN}✓ Test ticket found in worktree list${NC}"
else
  echo -e "${RED}✗ Test ticket not found in worktree list${NC}"
  exit 1
fi

# 3. Delete worktree
echo -e "\n${YELLOW}=== Testing worktree deletion ===${NC}"
"$BINARY_PATH" delete "$TEST_TICKET" -d

# Verify worktree no longer exists
if [ ! -d "$WORKTREE_PATH" ]; then
  echo -e "${GREEN}✓ Worktree directory deleted successfully${NC}"
else
  echo -e "${RED}✗ Worktree directory still exists at $WORKTREE_PATH${NC}"
  exit 1
fi

# 4. List again to verify it's gone
WORKTREE_LIST=$("$BINARY_PATH" list)
if echo "$WORKTREE_LIST" | grep -q "$TEST_TICKET"; then
  echo -e "${RED}✗ Test ticket still found in worktree list${NC}"
  exit 1
else
  echo -e "${GREEN}✓ Test ticket no longer in worktree list${NC}"
fi

echo -e "\n${GREEN}All integration tests passed!${NC}" 