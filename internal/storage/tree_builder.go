package storage

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/AdeshDeshmukh/helix/internal/objects"
)

type TreeBuilder struct {
	db      *Database
	entries map[string]TreeBuilderEntry
}

type TreeBuilderEntry struct {
	Path string
	Hash string
	Mode string
}

func NewTreeBuilder(db *Database) *TreeBuilder {
	return &TreeBuilder{
		db:      db,
		entries: make(map[string]TreeBuilderEntry),
	}
}

func (tb *TreeBuilder) AddEntry(path, hash, mode string) {
	tb.entries[path] = TreeBuilderEntry{
		Path: path,
		Hash: hash,
		Mode: mode,
	}
}

func (tb *TreeBuilder) BuildTree() (*objects.Tree, error) {
	tree := objects.NewTree()

	for _, entry := range tb.entries {
		tree.AddEntry(entry.Mode, filepath.Base(entry.Path), entry.Hash, "blob")
	}

	tree.Format()
	return tree, nil
}

func (tb *TreeBuilder) BuildTreeFromDirectory(dirPath string) (*objects.Tree, error) {
	tree := objects.NewTree()

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.Name() == ".helix" {
			continue
		}

		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())

		if entry.IsDir() {
			subtree, err := tb.BuildTreeFromDirectory(fullPath)
			if err != nil {
				return nil, fmt.Errorf("failed to build subtree for %s: %w", entry.Name(), err)
			}

			if err := tb.db.WriteTree(subtree); err != nil {
				return nil, fmt.Errorf("failed to write subtree: %w", err)
			}

			tree.AddEntry("040000", entry.Name(), subtree.Hash, "tree")
		} else {
			blob, err := objects.NewBlobFromFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("failed to create blob for %s: %w", entry.Name(), err)
			}

			if err := tb.db.WriteBlob(blob); err != nil {
				return nil, fmt.Errorf("failed to write blob: %w", err)
			}

			info, err := entry.Info()
			if err != nil {
				return nil, fmt.Errorf("failed to get file info: %w", err)
			}

			mode := GetFileMode(info)
			tree.AddEntry(mode, entry.Name(), blob.Hash, "blob")
		}
	}

	tree.Format()
	return tree, nil
}

func GetFileMode(info fs.FileInfo) string {
	mode := info.Mode()

	if mode.IsDir() {
		return "040000"
	}

	if mode&0111 != 0 {
		return "100755"
	}

	return "100644"
}
