package pull

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	utils "github.com/teaglebuilt/gh-subrepo/internal"
)

func TestPullCmd_NoGitRepo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-pull-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	cwd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(cwd)

	assert.NoError(t, os.Chdir(tmpDir))

	_, err = utils.GitRepoRoot()
	assert.Error(t, err, "Expected error when not in a git repo")
}

func TestPullCmd_MissingGitrepoFile(t *testing.T) {
	repoRoot, err := os.MkdirTemp("", "test-pull-gitrepo-*")
	assert.NoError(t, err)
	defer os.RemoveAll(repoRoot)

	subdir := "vendor/lib"
	assert.NoError(t, os.MkdirAll(filepath.Join(repoRoot, subdir), 0o755))

	err = utils.ExecCmd(repoRoot, "git", "init", "-b", "main")
	assert.NoError(t, err)

	gitrepoPath := filepath.Join(repoRoot, subdir, ".gitrepo")
	_, err = os.Stat(gitrepoPath)
	assert.True(t, os.IsNotExist(err), "Expected .gitrepo file not to exist")
}
