package clone

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	utils "github.com/teaglebuilt/gh-subrepo/internal"
)

func CloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone [repo-url] [subdir]",
		Short: "Clone a subrepo into a subdirectory",
		Args:  cobra.RangeArgs(1, 2),
		Run:   cloneRepository,
	}
	return cmd
}

func cloneRepository(cmd *cobra.Command, args []string) {
	repoURL := args[0]
	var subdir string

	if len(args) == 2 {
		subdir = args[1]
	} else {
		parts := strings.Split(repoURL, "/")
		lastPart := parts[len(parts)-1]
		subdir = strings.TrimSuffix(lastPart, ".git")
	}

	remote, _ := cmd.Flags().GetBool("remote")
	if remote {
		fmt.Println("Remote mode is not yet implemented.")
		return
	}

	fmt.Printf("Cloning subrepo from %s into %s...\n", repoURL, subdir)

	if _, err := os.Stat(subdir); err == nil {
		fmt.Printf("Subdirectory %s already exists. Abort.\n", subdir)
		os.Exit(1)
	}

	tmpDir := filepath.Join(os.TempDir(), "gh-subrepo-clone")
	defer os.RemoveAll(tmpDir)

	fmt.Println("Fetching repository temporarily...")
	if err := utils.ExecCmd("git", "clone", "--depth=1", repoURL, tmpDir); err != nil {
		fmt.Printf("Failed to clone temporary repository: %v\n", err)
		os.Exit(1)
	}

	if err := os.RemoveAll(filepath.Join(tmpDir, ".git")); err != nil {
		fmt.Printf("Failed to clean up temp repository: %v\n", err)
		os.Exit(1)
	}

	if err := utils.ExecCmd("mv", tmpDir, subdir); err != nil {
		fmt.Printf("Failed to move subrepo to subdirectory: %v\n", err)
		os.Exit(1)
	}

	gitRepoContent := fmt.Sprintf("[subrepo]\nremote = %s\nbranch = main\n", repoURL)
	gitRepoPath := filepath.Join(subdir, ".gitrepo")
	if err := os.WriteFile(gitRepoPath, []byte(gitRepoContent), 0o644); err != nil {
		fmt.Printf("Failed to write .gitrepo file: %v\n", err)
		os.Exit(1)
	}

	if err := utils.ExecCmd("git", "add", subdir); err != nil {
		fmt.Printf("Failed to add subdir to git: %v\n", err)
		os.Exit(1)
	}

	commitMsg := fmt.Sprintf("Add subrepo %s into %s", repoURL, subdir)
	if err := utils.ExecCmd("git", "commit", "-m", commitMsg); err != nil {
		fmt.Printf("Failed to commit subrepo: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Subrepo cloned and committed successfully.")
}
