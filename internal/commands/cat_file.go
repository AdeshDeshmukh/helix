package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/AdeshDeshmukh/helix/internal/storage"
	"github.com/spf13/cobra"
)

var catFileCmd = &cobra.Command{
	Use:   "cat-file <object-hash>",
	Short: "Display object content",
	Args:  cobra.ExactArgs(1),
	RunE:  runCatFile,
}

var catFilePretty bool

func init() {
	catFileCmd.Flags().BoolVarP(&catFilePretty, "pretty-print", "p", false, "Pretty-print object")
}

func getObjectType(data []byte) string {
	for i, b := range data {
		if b == ' ' {
			return string(data[:i])
		}
	}
	return ""
}

func runCatFile(cmd *cobra.Command, args []string) error {
	hash := args[0]

	repoPath, err := findRepositoryRoot(".")
	if err != nil {
		return fmt.Errorf("not a helix repository")
	}

	db := storage.NewDatabase(repoPath)

	if !db.ObjectExists(hash) {
		return fmt.Errorf("object %s not found", hash)
	}

	objectPath := storage.ObjectPath(repoPath, hash)
	file, err := os.Open(objectPath)
	if err != nil {
		return fmt.Errorf("failed to open object: %w", err)
	}
	defer file.Close()

	reader, err := storage.NewZlibReader(file)
	if err != nil {
		return fmt.Errorf("failed to create reader: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read object: %w", err)
	}

	objType := getObjectType(data)

	switch objType {
	case "blob":
		blob, err := db.ReadBlob(hash)
		if err != nil {
			return fmt.Errorf("failed to read blob: %w", err)
		}
		if catFilePretty {
			fmt.Printf("blob %d\n\n", blob.Size)
		}
		fmt.Print(string(blob.Content))
		return nil

	case "tree":
		tree, err := db.ReadTree(hash)
		if err != nil {
			return fmt.Errorf("failed to read tree: %w", err)
		}
		if catFilePretty {
			for _, entry := range tree.Entries {
				fmt.Printf("%s %s %s\t%s\n", entry.Mode, entry.Type, entry.Hash, entry.Name)
			}
		}
		return nil

	case "commit":
		commit, err := db.ReadCommit(hash)
		if err != nil {
			return fmt.Errorf("failed to read commit: %w", err)
		}
		if catFilePretty {
			fmt.Printf("tree %s\n", commit.Tree)
			if commit.Parent != "" {
				fmt.Printf("parent %s\n", commit.Parent)
			}
			fmt.Printf("author %s\n", commit.Author)
			fmt.Printf("committer %s\n", commit.Committer)
			fmt.Printf("\n%s\n", commit.Message)
		}
		return nil

	default:
		return fmt.Errorf("unknown object type: %s", objType)
	}
}

