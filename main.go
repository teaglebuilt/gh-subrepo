package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:                "subrepo",
		Short:              "GH extension wrapper for git-subrepo",
		DisableFlagParsing: true, // allow unknown flags
		Args:               cobra.ArbitraryArgs,
		TraverseChildren:   true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return delegateToGitSubrepo(args)
		},
	}

	rootCmd.AddCommand(initCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Install git-subrepo locally for GH extension",
		RunE: func(cmd *cobra.Command, args []string) error {
			installPath := filepath.Join(".", "git-subrepo")

			if _, err := os.Stat(installPath); err == nil {
				fmt.Println("git-subrepo already installed.")
				return nil
			}

			fmt.Println("Cloning git-subrepo...")
			if err := exec.Command("git", "clone", "https://github.com/ingydotnet/git-subrepo", installPath).Run(); err != nil {
				return fmt.Errorf("failed to clone git-subrepo: %w", err)
			}

			fmt.Println("git-subrepo installed to:", installPath)
			fmt.Println("You can now run `gh subrepo clone ...` etc.")
			fmt.Println("reload your shell if it is not registered on system path")
			return nil
		},
	}
}

func delegateToGitSubrepo(args []string) error {
	cwd, _ := os.Getwd()
	rcPath := filepath.Join(cwd, "git-subrepo", ".rc")

	// Ensure git-subrepo is installed
	if _, err := os.Stat(rcPath); os.IsNotExist(err) {
		return fmt.Errorf("git-subrepo is not installed. Run `gh subrepo init` first.")
	}

	// Build command to source the .rc and run git subrepo ...
	quotedArgs := shellQuoteArgs(args)
	command := fmt.Sprintf("source %q && git subrepo %s", rcPath, strings.Join(quotedArgs, " "))

	// Execute in a shell
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func shellQuoteArgs(args []string) []string {
	var out []string
	for _, arg := range args {
		out = append(out, fmt.Sprintf("'%s'", strings.ReplaceAll(arg, `'`, `'\''`)))
	}
	return out
}
