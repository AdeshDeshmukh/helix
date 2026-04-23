package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AdeshDeshmukh/helix/internal/objects"
	"github.com/AdeshDeshmukh/helix/internal/storage"
	"github.com/spf13/cobra"
)

var hashObjectCmd = &cobra.Command{
	Use:   "hash-object [file]",
	Short: "Compute object hash and optionally write to database",
	RunE:  runHashObject,
}

var hashObjectWrite bool
var hashObjectType string

func init() {
	hashObjectCmd.Flags().BoolVarP(&hashObjectWrite, "write", "w", false, "Write the object to database")
	hashObjectCmd.Flags().StringVarP(&hashObjectType, "type", "t", "blob", "Specify the type of object")
	rootCmd.AddCommand(hashObjectCmd)
}

func findRepositoryRoot(startPath string) (string, error) {
	currentPath, err := filepath.Abs(startPath)
	if err != nil {
		return "", err
	}

	for {
		helixPath := filepath.Join(currentPath, ".helix")
		if info, err := os.Stat(helixPath); err == nil && info.IsDir() {
			return currentPath, nil
		}

		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			return "", fmt.Errorf("not a helix repository (or any parent)")
		}
		currentPath = parentPath
	}
}

func runHashObject(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("file path required")
	}

	filePath := args[0]

	var repoPath string
	var db *storage.Database

	if hashObjectWrite {
		var err error
		repoPath, err = findRepositoryRoot(".")
		if err != nil {
			return fmt.Errorf("not a helix repository")
		}
		db = storage.NewDatabase(repoPath)
	}

	if hashObjectType == "blob" {
		blob, err := objects.NewBlobFromFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to create blob: %w", err)
		}

		if hashObjectWrite {
			if err := db.WriteBlob(blob); err != nil {
				return fmt.Errorf("failed to write blob: %w", err)
			}
		}

		fmt.Println(blob.Hash)
	} else {
		return fmt.Errorf("unsupported object type: %s", hashObjectType)
	}

	return nil
}
