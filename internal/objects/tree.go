package objects

import (
	"fmt"
	"sort"
	"strings"

	"github.com/AdeshDeshmukh/helix/internal/utils"
)

type TreeEntry struct {
	Mode string
	Name string
	Hash string
	Type string
}

type Tree struct {
	Entries []TreeEntry
	Hash    string
	Size    int
}

func NewTree() *Tree {
	return &Tree{
		Entries: make([]TreeEntry, 0),
	}
}

func (t *Tree) AddEntry(mode, name, hash, objType string) {
	entry := TreeEntry{
		Mode: mode,
		Name: name,
		Hash: hash,
		Type: objType,
	}
	t.Entries = append(t.Entries, entry)
}

func (t *Tree) Sort() {
	sort.Slice(t.Entries, func(i, j int) bool {
		nameA := t.Entries[i].Name
		nameB := t.Entries[j].Name
		
		if t.Entries[i].Mode == "040000" {
			nameA += "/"
		}
		if t.Entries[j].Mode == "040000" {
			nameB += "/"
		}
		
		return nameA < nameB
	})
}

func (t *Tree) Type() string {
	return "tree"
}

func (t *Tree) Format() []byte {
	var content []byte
	
	t.Sort()
	
	for _, entry := range t.Entries {
		entryHeader := fmt.Sprintf("%s %s\x00", entry.Mode, entry.Name)
		content = append(content, []byte(entryHeader)...)
		
		hashBytes, err := utils.HexToBytes(entry.Hash)
		if err != nil {
			continue
		}
		content = append(content, hashBytes...)
	}
	
	t.Size = len(content)
	t.Hash = utils.HashContent("tree", content)
	
	header := fmt.Sprintf("tree %d\x00", t.Size)
	return append([]byte(header), content...)
}

func (t *Tree) String() string {
	var entries []string
	for _, entry := range t.Entries {
		entries = append(entries, fmt.Sprintf("%s %s %s %s", 
			entry.Mode, entry.Type, entry.Hash[:8], entry.Name))
	}
	return fmt.Sprintf("Tree{hash: %s, entries: %d}\n%s", 
		t.Hash, len(t.Entries), strings.Join(entries, "\n"))
}