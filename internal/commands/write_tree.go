package commands

import (
	"fmt"

	"github.com/AdeshDeshmukh/helix/internal/storage"
	"github.com/spf13/cobra"
)

var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "Create tree object from current directory",
	RunE:  runWriteTree,
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}

func runWriteTree(cmd *cobra.Command, args []string) error {
	repoPath, err := findRepositoryRoot(".")
	if err != nil {
		return fmt.Errorf("not a helix repository")
	}

	db := storage.NewDatabase(repoPath)
	builder := storage.NewTreeBuilder(db)

	tree, err := builder.BuildTreeFromDirectory(".")
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err)
	}

	if err := db.WriteTree(tree); err != nil {
		return fmt.Errorf("failed to write tree: %w", err)
	}

	fmt.Println(tree.Hash)
	return nil
}