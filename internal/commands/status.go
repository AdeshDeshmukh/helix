package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/AdeshDeshmukh/helix/internal/index"
	"github.com/AdeshDeshmukh/helix/internal/objects"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		repoPath, err := findRepositoryRoot(wd)
		if err != nil {
			return err
		}

		idx := index.NewIndex(repoPath)
		err = idx.Load()
		if err != nil {
			return err
		}

		stagedFiles := make(map[string]bool)
		for _, entry := range idx.GetEntries() {
			stagedFiles[entry.Path] = true
		}

		allFiles := make(map[string]bool)
		err = filepath.WalkDir(wd, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				if d.Name() == ".helix" || strings.HasPrefix(d.Name(), ".") {
					return filepath.SkipDir
				}
				return nil
			}

			if strings.HasPrefix(d.Name(), ".") {
				return nil
			}

			relPath, _ := filepath.Rel(wd, path)
			allFiles[relPath] = true
			return nil
		})
		if err != nil {
			return err
		}

		changedToCommit := []string{}
		changedNotStaged := []string{}
		untracked := []string{}

		for file := range allFiles {
			fullPath := filepath.Join(wd, file)

			if stagedFiles[file] {
				blob, err := objects.NewBlobFromFile(fullPath)
				if err == nil {
					entry, _ := idx.GetEntry(file)
					if entry.Hash != blob.Hash {
						changedNotStaged = append(changedNotStaged, file)
					} else {
						changedToCommit = append(changedToCommit, file)
					}
				}
			} else {
				untracked = append(untracked, file)
			}
		}

		if len(changedToCommit) > 0 {
			fmt.Println("Changes to be committed:")
			for _, file := range changedToCommit {
				fmt.Printf("  \033[32m%s\033[0m\n", file)
			}
			fmt.Println()
		}

		if len(changedNotStaged) > 0 {
			fmt.Println("Changes not staged for commit:")
			for _, file := range changedNotStaged {
				fmt.Printf("  \033[33m%s\033[0m\n", file)
			}
			fmt.Println()
		}

		if len(untracked) > 0 {
			fmt.Println("Untracked files:")
			fmt.Println("  (use \"helix add <file>\" to include in commit)")
			for _, file := range untracked {
				fmt.Printf("  \033[31m%s\033[0m\n", file)
			}
		}

		if len(changedToCommit) == 0 && len(changedNotStaged) == 0 && len(untracked) == 0 {
			fmt.Println("nothing to commit, working tree clean")
		}

		return nil
	},
}
