package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push [subdir]",
	Short: "Push local changes back to the subrepo's remote",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subdir := args[0]

		remote, _ := cmd.Flags().GetBool("remote")
		if remote {
			fmt.Println("Remote mode is not yet implemented.")
			return
		}

		fmt.Printf("Pushing changes from subrepo %s...\n", subdir)

		gitCmd := exec.Command("git", "subrepo", "push", subdir)
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr

		if err := gitCmd.Run(); err != nil {
			fmt.Printf("Failed to push subrepo: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Subrepo pushed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
