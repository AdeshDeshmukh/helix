package commands

import (
	"fmt"

	"github.com/AdeshDeshmukh/helix/internal/storage"
	"github.com/spf13/cobra"
)

var lsTreeCmd = &cobra.Command{
	Use:   "ls-tree <tree-hash>",
	Short: "List the contents of a tree object",
	Args:  cobra.ExactArgs(1),
	RunE:  runLsTree,
}

func init() {
	rootCmd.AddCommand(lsTreeCmd)
}

func runLsTree(cmd *cobra.Command, args []string) error {
	treeHash := args[0]

	repoPath, err := findRepositoryRoot(".")
	if err != nil {
		return fmt.Errorf("not a helix repository")
	}

	db := storage.NewDatabase(repoPath)

	tree, err := db.ReadTree(treeHash)
	if err != nil {
		return fmt.Errorf("failed to read tree: %w", err)
	}

	for _, entry := range tree.Entries {
		fmt.Printf("%s %s %s\t%s\n",
			entry.Mode, entry.Type, entry.Hash, entry.Name)
	}

	return nil
}