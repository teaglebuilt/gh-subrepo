package status

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	utils "github.com/teaglebuilt/gh-subrepo/internal"
)

func TestFindAllSubrepos(t *testing.T) {
	rootDir, _ := os.MkdirTemp("", "test-status-*")
	defer os.RemoveAll(rootDir)

	os.MkdirAll(filepath.Join(rootDir, "vendor/lib1"), 0755)
	os.WriteFile(filepath.Join(rootDir, "vendor/lib1/.gitrepo"), []byte("remote = url"), 0644)

	os.MkdirAll(filepath.Join(rootDir, "vendor/lib2"), 0755)
	os.WriteFile(filepath.Join(rootDir, "vendor/lib2/.gitrepo"), []byte("remote = url"), 0644)

	subdirs := findAllSubrepos(rootDir)
	assert.Len(t, subdirs, 2)
	assert.Contains(t, subdirs, filepath.Join(rootDir, "vendor/lib1"))
	assert.Contains(t, subdirs, filepath.Join(rootDir, "vendor/lib2"))
}

func setupGitConfig(t *testing.T, repoDir string) {
	assert.NoError(t, utils.ExecCmd(repoDir, "git", "config", "user.name", "testuser"))
	assert.NoError(t, utils.ExecCmd(repoDir, "git", "config", "user.email", "testuser@example.com"))
}

func TestStatusCmd_Integration(t *testing.T) {
	subrepoDir := "vendor/lib"

	rootDir, err := os.MkdirTemp("", "test-status-integration-*")
	assert.NoError(t, err)
	defer os.RemoveAll(rootDir)

	assert.NoError(t, utils.ExecCmd(rootDir, "git", "init", "-b", "main"))
	setupGitConfig(t, rootDir)

	remoteRepo, err := os.MkdirTemp("", "remote-subrepo-*")
	assert.NoError(t, err)
	defer os.RemoveAll(remoteRepo)

	assert.NoError(t, utils.ExecCmd(remoteRepo, "git", "init", "--bare"))

	tmpWork, err := os.MkdirTemp("", "remote-workspace-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpWork)

	assert.NoError(t, utils.ExecCmd("", "git", "clone", remoteRepo, tmpWork))
	setupGitConfig(t, tmpWork)

	os.WriteFile(filepath.Join(tmpWork, "README.md"), []byte("Initial remote content"), 0644)
	assert.NoError(t, utils.ExecCmd(tmpWork, "git", "add", "."))
	assert.NoError(t, utils.ExecCmd(tmpWork, "git", "commit", "-m", "Initial commit"))
	assert.NoError(t, utils.ExecCmd(tmpWork, "git", "push", "-u", "origin", "main"))

	assert.NoError(t, utils.ExecCmd(rootDir, "git", "clone", "-b", "main", remoteRepo, filepath.Join(rootDir, subrepoDir)))

	os.RemoveAll(filepath.Join(rootDir, subrepoDir, ".git"))

	gitrepoContent := "[subrepo]\nremote = " + remoteRepo + "\nbranch = main\n"
	os.WriteFile(filepath.Join(rootDir, subrepoDir, ".gitrepo"), []byte(gitrepoContent), 0644)

	assert.NoError(t, utils.ExecCmd(rootDir, "git", "add", "."))
	assert.NoError(t, utils.ExecCmd(rootDir, "git", "commit", "-m", "Add subrepo"))

	checkStatus(rootDir, filepath.Join(rootDir, subrepoDir))

	os.WriteFile(filepath.Join(tmpWork, "CHANGELOG.md"), []byte("Changelog added"), 0644)
	assert.NoError(t, utils.ExecCmd(tmpWork, "git", "add", "."))
	assert.NoError(t, utils.ExecCmd(tmpWork, "git", "commit", "-m", "Add CHANGELOG"))
	assert.NoError(t, utils.ExecCmd(tmpWork, "git", "push", "origin", "main"))

	checkStatus(rootDir, filepath.Join(rootDir, subrepoDir))
}
