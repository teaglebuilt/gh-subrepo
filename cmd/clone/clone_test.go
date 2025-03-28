package clone

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	utils "github.com/teaglebuilt/gh-subrepo/internal"
)

func TestGitRepoRoot(t *testing.T) {
	root, err := utils.GitRepoRoot()
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
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			fmt.Printf("Warning: failed to remove tmp dir: %v\n", err)
		}
	}()

	testFile := filepath.Join(tmpDir, "testfile.txt")
	err = utils.ExecCmd(tmpDir, "touch", testFile)
	if err != nil {
		t.Fatalf("ExecCmd failed: %v", err)
	}

	if _, err := os.Stat(testFile); err != nil {
		t.Fatalf("File was not created: %v", err)
	}
}
