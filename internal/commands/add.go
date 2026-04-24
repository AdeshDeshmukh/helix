package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/AdeshDeshmukh/helix/internal/index"
	"github.com/AdeshDeshmukh/helix/internal/objects"
	"github.com/AdeshDeshmukh/helix/internal/storage"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("please specify files or directories to add")
		}

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		repoPath, err := findRepositoryRoot(wd)
		if err != nil {
			return err
		}

		db := storage.NewDatabase(repoPath)
		idx := index.NewIndex(repoPath)
		err = idx.Load()
		if err != nil {
			return err
		}

		filesToAdd := make(map[string]bool)

		for _, arg := range args {
			if arg == "." {
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
					filesToAdd[relPath] = true
					return nil
				})
				if err != nil {
					return err
				}
			} else {
				fullPath := arg
				if !filepath.IsAbs(arg) {
					fullPath = filepath.Join(wd, arg)
				}

				info, err := os.Stat(fullPath)
				if err != nil {
					return err
				}

				if info.IsDir() {
					err = filepath.WalkDir(fullPath, func(path string, d fs.DirEntry, err error) error {
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
						filesToAdd[relPath] = true
						return nil
					})
					if err != nil {
						return err
					}
				} else {
					if !strings.HasPrefix(info.Name(), ".") {
						relPath, _ := filepath.Rel(wd, fullPath)
						filesToAdd[relPath] = true
					}
				}
			}
		}

		for relPath := range filesToAdd {
			fullPath := filepath.Join(wd, relPath)

			blob, err := objects.NewBlobFromFile(fullPath)
			if err != nil {
				return fmt.Errorf("failed to create blob for %s: %w", relPath, err)
			}

			err = db.WriteBlob(blob)
			if err != nil {
				return fmt.Errorf("failed to write blob for %s: %w", relPath, err)
			}

			info, err := os.Stat(fullPath)
			if err != nil {
				return err
			}

			mode := storage.GetFileMode(info)
			idx.Add(relPath, blob.Hash, mode, info.Size(), info.ModTime())
		}

		err = idx.Save()
		if err != nil {
			return err
		}

		fmt.Printf("Staged %d file(s)\n", len(filesToAdd))
		return nil
	},
}
