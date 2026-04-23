package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/AdeshDeshmukh/helix/internal/objects"
	"github.com/AdeshDeshmukh/helix/internal/storage"
	"github.com/spf13/cobra"
)

var commitTreeCmd = &cobra.Command{
	Use:   "commit-tree <tree-hash>",
	Short: "Create a commit object from a tree",
	Args:  cobra.ExactArgs(1),
	RunE:  runCommitTree,
}

var commitParent string
var commitMessage string

func init() {
	commitTreeCmd.Flags().StringVarP(&commitParent, "parent", "p", "", "Parent commit hash")
	commitTreeCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Commit message")
}

func runCommitTree(cmd *cobra.Command, args []string) error {
	if commitMessage == "" {
		return fmt.Errorf("message is required (-m)")
	}

	treeHash := args[0]

	repoPath, err := findRepositoryRoot(".")
	if err != nil {
		return fmt.Errorf("not a helix repository")
	}

	authorName := os.Getenv("GIT_AUTHOR_NAME")
	if authorName == "" {
		authorName = "Adesh Deshmukh"
	}

	authorEmail := os.Getenv("GIT_AUTHOR_EMAIL")
	if authorEmail == "" {
		authorEmail = "adeshkd123@gmail.com"
	}

	timestamp := time.Now().Unix()
	author := fmt.Sprintf("%s <%s> %d +0000", authorName, authorEmail, timestamp)
	committer := fmt.Sprintf("%s <%s> %d +0000", authorName, authorEmail, timestamp)

	commit := objects.NewCommit(treeHash, commitParent, author, committer, commitMessage)

	db := storage.NewDatabase(repoPath)
	if err := db.WriteCommit(commit); err != nil {
		return fmt.Errorf("failed to write commit: %w", err)
	}

	fmt.Println(commit.Hash)
	return nil
}
