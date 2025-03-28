package status

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	utils "github.com/teaglebuilt/gh-subrepo/internal"
)

func setupGitConfig(t *testing.T, repoDir string) {
	assert.NoError(t, utils.ExecCmd(repoDir, "git", "config", "user.name", "testuser"))
	assert.NoError(t, utils.ExecCmd(repoDir, "git", "config", "user.email", "testuser@example.com"))
}

func TestStatusCmd_Integration(t *testing.T) {
	rootDir, err := os.MkdirTemp("", "test-status-integration-*")
	assert.NoError(t, err)
	defer os.RemoveAll(rootDir)

	assert.NoError(t, utils.ExecCmd(rootDir, "git", "init", "-b", "main"))
	setupGitConfig(t, rootDir)

	remoteRepo, err := os.MkdirTemp("", "remote-subrepo-*.git")
	assert.NoError(t, err)
	defer os.RemoveAll(remoteRepo)

	assert.NoError(t, utils.ExecCmd(remoteRepo, "git", "init", "--bare", "--initial-branch=main"))

	tmpWork, err := os.MkdirTemp("", "remote-workspace-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpWork)

	assert.NoError(t, utils.ExecCmd("", "git", "clone", remoteRepo, tmpWork))
	setupGitConfig(t, tmpWork)

	readmePath := filepath.Join(tmpWork, "README.md")
	assert.NoError(t, os.WriteFile(readmePath, []byte("Initial remote content"), 0644))

	assert.NoError(t, utils.ExecCmd(tmpWork, "git", "add", "."))
	assert.NoError(t, utils.ExecCmd(tmpWork, "git", "commit", "-m", "Initial commit"))
	assert.NoError(t, utils.ExecCmd(tmpWork, "git", "push", "-u", "origin", "main"))

	subrepoDir := "vendor/lib"
	subrepoPath := filepath.Join(rootDir, subrepoDir)
	assert.NoError(t, utils.ExecCmd(rootDir, "git", "clone", "-b", "main", remoteRepo, subrepoPath))

	os.RemoveAll(filepath.Join(subrepoPath, ".git"))
	gitrepoContent := "[subrepo]\nremote = " + remoteRepo + "\nbranch = main\n"
	os.WriteFile(filepath.Join(subrepoPath, ".gitrepo"), []byte(gitrepoContent), 0644)

	assert.NoError(t, utils.ExecCmd(rootDir, "git", "add", "."))
	assert.NoError(t, utils.ExecCmd(rootDir, "git", "commit", "-m", "Add subrepo"))

	checkStatus(subrepoPath)

	changelogPath := filepath.Join(tmpWork, "CHANGELOG.md")
	os.WriteFile(changelogPath, []byte("Changelog added"), 0644)
	assert.NoError(t, utils.ExecCmd(tmpWork, "git", "add", "."))
	assert.NoError(t, utils.ExecCmd(tmpWork, "git", "commit", "-m", "Add CHANGELOG"))
	assert.NoError(t, utils.ExecCmd(tmpWork, "git", "push", "origin", "main"))

	checkStatus(subrepoPath)
}
