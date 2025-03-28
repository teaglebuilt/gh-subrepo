package status

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	utils "github.com/teaglebuilt/gh-subrepo/internal"
	"gopkg.in/ini.v1"
)

func StatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [subdir]",
		Short: "Show status of subrepo(s)",
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			repoRoot, err := utils.GitRepoRoot()
			if err != nil {
				fmt.Printf("Not inside a git repository: %v\n", err)
				os.Exit(1)
			}

			var subdirs []string
			if len(args) == 1 {
				subdirs = []string{args[0]}
			} else {
				subdirs = findAllSubrepos(repoRoot)
			}

			for _, subdir := range subdirs {
				checkStatus(repoRoot, subdir)
			}
		},
	}
	return cmd
}

func findAllSubrepos(root string) []string {
	var dirs []string
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if filepath.Base(path) == ".gitrepo" {
			dirs = append(dirs, filepath.Dir(path))
		}
		return nil
	})
	return dirs
}

func checkStatus(repoRoot, subdir string) {
	gitrepoFile := filepath.Join(subdir, ".gitrepo")

	cfg, err := ini.Load(gitrepoFile)
	if err != nil {
		fmt.Printf("[%s] failed to load .gitrepo: %v\n", subdir, err)
		return
	}

	remote := cfg.Section("subrepo").Key("remote").String()
	branch := cfg.Section("subrepo").Key("branch").MustString("main")

	tmpDir, _ := os.MkdirTemp("", "subrepo-status-*")
	defer os.RemoveAll(tmpDir)

	if err := utils.ExecCmd("", "git", "clone", "--depth=1", "-b", branch, remote, tmpDir); err != nil {
		fmt.Printf("[%s] failed to clone remote: %v\n", subdir, err)
		return
	}

	localHash, _ := GitHash(subdir)
	remoteHash, _ := GitHash(tmpDir)

	if localHash == remoteHash {
		fmt.Printf("[%s] ✅ Up to date\n", subdir)
	} else {
		fmt.Printf("[%s] ⚠️ Out of sync\n", subdir)
	}
}

func GitHash(dir string) (string, error) {
	out, err := execCommandOutput(dir, "git", "rev-parse", "HEAD")
	return out, err
}

func execCommandOutput(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output[:40]), nil // 40-char hash
}
