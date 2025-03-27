package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

var pullCmd = &cobra.Command{
	Use:   "pull [subdir]",
	Short: "Pull latest changes into a subrepo",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subdir := args[0]

		repoRoot, err := GitRepoRoot()
		if err != nil {
			fmt.Printf("Not inside a git repository: %v\n", err)
			os.Exit(1)
		}

		subrepoPath := filepath.Join(repoRoot, subdir)
		gitrepoFile := filepath.Join(subrepoPath, ".gitrepo")

		cfg, err := ini.Load(gitrepoFile)
		if err != nil {
			fmt.Printf("Failed to load .gitrepo file: %v\n", err)
			os.Exit(1)
		}

		remote := cfg.Section("subrepo").Key("remote").String()
		branch := cfg.Section("subrepo").Key("branch").MustString("main")

		tmpDir, err := os.MkdirTemp("", "gh-subrepo-pull-*")
		if err != nil {
			fmt.Printf("Failed to create temp dir: %v\n", err)
			os.Exit(1)
		}
		defer os.RemoveAll(tmpDir)

		fmt.Printf("Pulling latest changes from %s...\n", remote)
		if err := ExecCmd("", "git", "clone", "--depth=1", "-b", branch, remote, tmpDir); err != nil {
			fmt.Printf("Failed to clone subrepo temporarily: %v\n", err)
			os.Exit(1)
		}

		if err := os.RemoveAll(filepath.Join(tmpDir, ".git")); err != nil {
			fmt.Printf("Failed to remove temp .git: %v\n", err)
			os.Exit(1)
		}

		if err := ExecCmd("", "rsync", "-av", "--delete", tmpDir+"/", subrepoPath+"/"); err != nil {
			fmt.Printf("Failed to sync files: %v\n", err)
			os.Exit(1)
		}

		if err := ExecCmd(repoRoot, "git", "add", subdir); err != nil {
			fmt.Printf("Failed to stage updated subrepo: %v\n", err)
			os.Exit(1)
		}

		commitMsg := fmt.Sprintf("Update subrepo %s from %s", subdir, remote)
		if err := ExecCmd(repoRoot, "git", "commit", "-m", commitMsg); err != nil {
			fmt.Printf("Failed to commit subrepo update: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Subrepo updated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
