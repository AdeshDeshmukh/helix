package index

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type IndexEntry struct {
	Path    string    `json:"path"`
	Hash    string    `json:"hash"`
	Mode    string    `json:"mode"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}

type Index struct {
	Entries map[string]IndexEntry
	RepoPath string
}

func NewIndex(repoPath string) *Index {
	return &Index{
		Entries:  make(map[string]IndexEntry),
		RepoPath: repoPath,
	}
}

func (idx *Index) Load() error {
	indexPath := filepath.Join(idx.RepoPath, ".helix", "index")
	
	data, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			idx.Entries = make(map[string]IndexEntry)
			return nil
		}
		return err
	}
	
	var entries map[string]IndexEntry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return fmt.Errorf("failed to parse index: %w", err)
	}
	
	idx.Entries = entries
	return nil
}

func (idx *Index) Save() error {
	indexPath := filepath.Join(idx.RepoPath, ".helix", "index")
	
	data, err := json.MarshalIndent(idx.Entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}
	
	err = os.WriteFile(indexPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write index: %w", err)
	}
	
	return nil
}

func (idx *Index) Add(path, hash, mode string, size int64, modTime time.Time) {
	idx.Entries[path] = IndexEntry{
		Path:    path,
		Hash:    hash,
		Mode:    mode,
		Size:    size,
		ModTime: modTime,
	}
}

func (idx *Index) Remove(path string) {
	delete(idx.Entries, path)
}

func (idx *Index) GetEntry(path string) (IndexEntry, bool) {
	entry, exists := idx.Entries[path]
	return entry, exists
}

func (idx *Index) GetEntries() []IndexEntry {
	var entries []IndexEntry
	for _, entry := range idx.Entries {
		entries = append(entries, entry)
	}
	return entries
}

func (idx *Index) Clear() {
	idx.Entries = make(map[string]IndexEntry)
}
