package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:                "subrepo",
		Short:              "Wraps git-subrepo for GitHub CLI",
		Long:               "A gh extension to delegate subcommands to git-subrepo, installing it if necessary.",
		DisableFlagParsing: true, // allow unknown flags through to git-subrepo
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Fprintln(os.Stderr, "No command provided. Usage: gh subrepo <command> [args]")
				os.Exit(1)
			}

			if err := execGitSubrepo(args); err != nil {
				fmt.Fprintf(os.Stderr, "Error running git-subrepo: %v\n", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.AddCommand(InitCommand())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func InitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Install and configure git-subrepo locally",
		RunE: func(cmd *cobra.Command, args []string) error {
			return installGitSubrepo()
		},
	}
}

func execGitSubrepo(args []string) error {
	if !isGitSubrepoInstalled() {
		return fmt.Errorf("failed to install git-subrepo. Run init command first!: %w")
	}

	// Set PATH to include our fallback
	binPath := filepath.Join(os.Getenv("HOME"), ".gh-subrepo", "bin")
	pathEnv := fmt.Sprintf("%s:%s", binPath, os.Getenv("PATH"))
	cmd := exec.Command("git", append([]string{"subrepo"}, args...)...)
	cmd.Env = append(os.Environ(), "PATH="+pathEnv)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func isGitSubrepoInstalled() bool {
	cmd := exec.Command("git", "subrepo", "--version")
	cmd.Stderr = nil
	cmd.Stdout = nil
	return cmd.Run() == nil
}

func installGitSubrepo() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	installPath := filepath.Join(cwd, "git-subrepo")

	// Clone the git-subrepo repo if not already installed
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", "https://github.com/ingydotnet/git-subrepo", installPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git clone failed: %w", err)
		}
	}

	// Add sourcing line to user's shell config
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", "https://github.com/ingydotnet/git-subrepo", installPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git clone failed: %w", err)
		}
	}

	// Manually source the .rc file just for this runtime
	// Add ./git-subrepo to PATH for this session
	gitSubrepoRc := filepath.Join(installPath, ".rc")
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && exec git subrepo --version", gitSubrepoRc))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
