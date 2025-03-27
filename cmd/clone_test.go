package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGitRepoRoot(t *testing.T) {
	root, err := GitRepoRoot()
	if err != nil {
		t.Fatalf("Failed to find git root: %v", err)
	}

	if root == "" {
		t.Fatal("Expected git root to not be empty")
	}
}

func TestExecCmd(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-gh-subrepo-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "testfile.txt")
	err = ExecCmd(tmpDir, "touch", testFile)
	if err != nil {
		t.Fatalf("ExecCmd failed: %v", err)
	}

	if _, err := os.Stat(testFile); err != nil {
		t.Fatalf("File was not created: %v", err)
	}
}
