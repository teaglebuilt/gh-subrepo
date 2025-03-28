package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/teaglebuilt/gh-subrepo/cmd/clone"
	"github.com/teaglebuilt/gh-subrepo/cmd/fetch"
	"github.com/teaglebuilt/gh-subrepo/cmd/status"
)

var rootCmd = &cobra.Command{
	Use:   "subrepo",
	Short: "Manage git subrepos with GitHub CLI",
	Long:  `subrepo provides easy subrepo management within GitHub CLI.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("remote", "r", false, "Use remote mode via GitHub API")
	rootCmd.AddCommand(
		clone.CloneCmd(),
		fetch.FetchCmd(),
		status.StatusCmd(),
	)
}
