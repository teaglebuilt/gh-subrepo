package fetch

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	utils "github.com/teaglebuilt/gh-subrepo/internal"
	"gopkg.in/ini.v1"
)

func FetchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fetch [subdir]",
		Short: "Fetch latest subrepo changes (without merging)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			subdir := args[0]

			repoRoot, err := utils.GitRepoRoot()
			if err != nil {
				fmt.Printf("Not inside a git repository: %v\n", err)
				os.Exit(1)
			}

			gitrepoFile := filepath.Join(repoRoot, subdir, ".gitrepo")

			cfg, err := ini.Load(gitrepoFile)
			if err != nil {
				fmt.Printf("Failed to load .gitrepo file: %v\n", err)
				os.Exit(1)
			}

			remote := cfg.Section("subrepo").Key("remote").String()
			branch := cfg.Section("subrepo").Key("branch").MustString("main")

			if err := utils.ExecCmd(repoRoot, "git", "fetch", remote, branch); err != nil {
				fmt.Printf("Failed to fetch subrepo: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Subrepo fetched successfully.")
		},
	}
	return cmd
}
