package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/ini.v1"
)

func TestPushGitRepoFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-push-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	subrepoDir := filepath.Join(tmpDir, "subrepo")
	os.Mkdir(subrepoDir, 0755)

	gitrepoContent := `[subrepo]
remote = https://github.com/example/pushrepo.git
branch = develop
`
	gitrepoPath := filepath.Join(subrepoDir, ".gitrepo")
	if err := os.WriteFile(gitrepoPath, []byte(gitrepoContent), 0644); err != nil {
		t.Fatalf("Failed to write .gitrepo: %v", err)
	}

	cfg, err := ini.Load(gitrepoPath)
	if err != nil {
		t.Fatalf("Failed to load .gitrepo: %v", err)
	}

	remote := cfg.Section("subrepo").Key("remote").String()
	if remote != "https://github.com/example/pushrepo.git" {
		t.Fatalf("Expected remote URL to match, got %s", remote)
	}

	branch := cfg.Section("subrepo").Key("branch").String()
	if branch != "develop" {
		t.Fatalf("Expected branch to be 'develop', got %s", branch)
	}
}
