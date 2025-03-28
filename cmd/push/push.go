package push

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	utils "github.com/teaglebuilt/gh-subrepo/internal"
	"gopkg.in/ini.v1"
)

func PushCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push [subdir]",
		Short: "Push local changes back to the subrepo's remote",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			subdir := args[0]

			repoRoot, err := utils.GitRepoRoot()
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

			tmpDir, err := os.MkdirTemp("", "gh-subrepo-push-*")
			if err != nil {
				fmt.Printf("Failed to create temp dir: %v\n", err)
				os.Exit(1)
			}
			defer os.RemoveAll(tmpDir)

			fmt.Printf("Preparing to push changes to %s...\n", remote)

			if err := utils.ExecCmd("", "git", "clone", "-b", branch, remote, tmpDir); err != nil {
				fmt.Printf("Failed to clone subrepo remotely: %v\n", err)
				os.Exit(1)
			}

			if err := utils.ExecCmd("", "rsync", "-av", "--delete", subrepoPath+"/", tmpDir+"/"); err != nil {
				fmt.Printf("Failed to sync local changes: %v\n", err)
				os.Exit(1)
			}

			if err := utils.ExecCmd(tmpDir, "git", "add", "."); err != nil {
				fmt.Printf("Failed to stage changes in temp repo: %v\n", err)
				os.Exit(1)
			}

			commitMsg := fmt.Sprintf("Sync updates from main repo for %s", subdir)
			if err := utils.ExecCmd(tmpDir, "git", "commit", "-m", commitMsg); err != nil {
				fmt.Printf("No changes to push or commit failed: %v\n", err)
				os.Exit(1)
			}

			if err := utils.ExecCmd(tmpDir, "git", "push", "origin", branch); err != nil {
				fmt.Printf("Failed to push changes: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Changes pushed to subrepo successfully.")
		},
	}
	return cmd
}
