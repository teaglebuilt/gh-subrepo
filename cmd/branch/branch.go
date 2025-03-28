package branch

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	utils "github.com/teaglebuilt/gh-subrepo/internal"
)

func BranchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branch [subdir] [branch-name]",
		Short: "Create a branch reflecting subrepo state",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			subdir := args[0]
			branchName := args[1]

			repoRoot, err := utils.GitRepoRoot()
			if err != nil {
				fmt.Printf("Not a git repository: %v\n", err)
				os.Exit(1)
			}

			fullSubdirPath := filepath.Join(repoRoot, subdir)

			if err := utils.ExecCmd(repoRoot, "git", "subtree", "split", "-P", subdir, "-b", branchName); err != nil {
				fmt.Printf("Failed to create branch from subdir: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Branch '%s' created from '%s'.\n", branchName, fullSubdirPath)
		},
	}
	return cmd
}
