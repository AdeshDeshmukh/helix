package commands

import (
	"fmt"

	"github.com/AdeshDeshmukh/helix/internal/refs"
	"github.com/AdeshDeshmukh/helix/internal/storage"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runLog,
}

func runLog(cmd *cobra.Command, args []string) error {
	var startHash string

	if len(args) > 0 {
		startHash = args[0]
	} else {
		repoPath, err := findRepositoryRoot(".")
		if err != nil {
			return fmt.Errorf("not a helix repository")
		}

		hash, err := refs.GetHEAD(repoPath)
		if err != nil {
			return fmt.Errorf("failed to read HEAD: %w", err)
		}

		if hash == "" {
			return fmt.Errorf("no commits yet")
		}

		startHash = hash
	}

	repoPath, err := findRepositoryRoot(".")
	if err != nil {
		return fmt.Errorf("not a helix repository")
	}

	db := storage.NewDatabase(repoPath)

	currentHash := startHash
	for currentHash != "" {
		commit, err := db.ReadCommit(currentHash)
		if err != nil {
			return fmt.Errorf("failed to read commit %s: %w", currentHash, err)
		}

		fmt.Printf("commit %s\n", currentHash)
		fmt.Printf("Author: %s\n", commit.Author)
		fmt.Printf("Committer: %s\n", commit.Committer)
		fmt.Printf("\n    %s\n\n", commit.Message)

		currentHash = commit.Parent
	}

	return nil
}
