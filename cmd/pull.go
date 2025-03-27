package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull [subdir]",
	Short: "Pull latest changes into a subrepo",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subdir := args[0]

		remote, _ := cmd.Flags().GetBool("remote")
		if remote {
			fmt.Println("Remote mode is not yet implemented.")
			return
		}

		fmt.Printf("Pulling latest changes into subrepo %s...\n", subdir)

		gitCmd := exec.Command("git", "subrepo", "pull", subdir)
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr

		if err := gitCmd.Run(); err != nil {
			fmt.Printf("Failed to pull subrepo: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Subrepo pulled successfully.")
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
