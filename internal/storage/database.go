
package storage

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/AdeshDeshmukh/helix/internal/objects"
	"github.com/AdeshDeshmukh/helix/internal/utils"
)

func ObjectPath(repoPath, hash string) string {
	return utils.ObjectPath(repoPath, hash)
}

func NewZlibReader(reader io.Reader) (io.ReadCloser, error) {
	return zlib.NewReader(reader)
}

type Database struct {
	RepoPath string
}

func NewDatabase(repoPath string) *Database {
	return &Database{RepoPath: repoPath}
}

func (db *Database) WriteBlob(blob *objects.Blob) error {
	objectPath := utils.ObjectPath(db.RepoPath, blob.Hash)

	if err := os.MkdirAll(filepath.Dir(objectPath), 0755); err != nil {
		return fmt.Errorf("failed to create object directory: %w", err)
	}

	if _, err := os.Stat(objectPath); err == nil {
		return nil
	}

	file, err := os.Create(objectPath)
	if err != nil {
		return fmt.Errorf("failed to create object file: %w", err)
	}
	defer file.Close()

	writer := zlib.NewWriter(file)
	defer writer.Close()

	if _, err := writer.Write(blob.Format()); err != nil {
		return fmt.Errorf("failed to write object content: %w", err)
	}

	return nil
}

func (db *Database) ReadBlob(hash string) (*objects.Blob, error) {
	if !utils.ValidateHash(hash) {
		return nil, fmt.Errorf("invalid hash: %s", hash)
	}

	objectPath := utils.ObjectPath(db.RepoPath, hash)

	file, err := os.Open(objectPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("object %s not found", hash)
		}
		return nil, fmt.Errorf("failed to open object file: %w", err)
	}
	defer file.Close()

	reader, err := zlib.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create zlib reader: %w", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read object content: %w", err)
	}

	nullIndex := -1
	for i, b := range content {
		if b == 0 {
			nullIndex = i
			break
		}
	}

	if nullIndex == -1 {
		return nil, fmt.Errorf("invalid blob format: missing null separator")
	}

	blobContent := content[nullIndex+1:]
	blob := objects.NewBlobFromContent(blobContent)

	return blob, nil
}

func (db *Database) WriteTree(tree *objects.Tree) error {
	objectPath := utils.ObjectPath(db.RepoPath, tree.Hash)

	if err := os.MkdirAll(filepath.Dir(objectPath), 0755); err != nil {
		return fmt.Errorf("failed to create object directory: %w", err)
	}

	if _, err := os.Stat(objectPath); err == nil {
		return nil
	}

	file, err := os.Create(objectPath)
	if err != nil {
		return fmt.Errorf("failed to create object file: %w", err)
	}
	defer file.Close()

	writer := zlib.NewWriter(file)
	defer writer.Close()

	if _, err := writer.Write(tree.Format()); err != nil {
		return fmt.Errorf("failed to write object content: %w", err)
	}

	return nil
}

func (db *Database) ReadTree(hash string) (*objects.Tree, error) {
	if !utils.ValidateHash(hash) {
		return nil, fmt.Errorf("invalid hash: %s", hash)
	}

	objectPath := utils.ObjectPath(db.RepoPath, hash)

	file, err := os.Open(objectPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("object %s not found", hash)
		}
		return nil, fmt.Errorf("failed to open object file: %w", err)
	}
	defer file.Close()

	reader, err := zlib.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create zlib reader: %w", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read object content: %w", err)
	}

	nullIndex := -1
	for i, b := range content {
		if b == 0 {
			nullIndex = i
			break
		}
	}

	if nullIndex == -1 {
		return nil, fmt.Errorf("invalid tree format: missing null separator")
	}

	tree := objects.NewTree()
	entryData := content[nullIndex+1:]

	for len(entryData) > 0 {
		spaceIndex := -1
		nullIndex := -1

		for i, b := range entryData {
			if b == ' ' && spaceIndex == -1 {
				spaceIndex = i
			}
			if b == 0 {
				nullIndex = i
				break
			}
		}

		if spaceIndex == -1 || nullIndex == -1 || nullIndex-spaceIndex <= 1 {
			break
		}

		mode := string(entryData[:spaceIndex])
		name := string(entryData[spaceIndex+1 : nullIndex])

		if len(entryData) < nullIndex+21 {
			break
		}

		hashBytes := entryData[nullIndex+1 : nullIndex+21]
		hash := utils.BytesToHex(hashBytes)

		objType := "blob"
		if mode == "040000" {
			objType = "tree"
		}

		tree.AddEntry(mode, name, hash, objType)
		entryData = entryData[nullIndex+21:]
	}

	return tree, nil
}

func (db *Database) ObjectExists(hash string) bool {
	if !utils.ValidateHash(hash) {
		return false
	}
	objectPath := utils.ObjectPath(db.RepoPath, hash)
	_, err := os.Stat(objectPath)
	return err == nil
}

func (db *Database) WriteCommit(commit *objects.Commit) error {
	objectPath := utils.ObjectPath(db.RepoPath, commit.Hash)

	if err := os.MkdirAll(filepath.Dir(objectPath), 0755); err != nil {
		return fmt.Errorf("failed to create object directory: %w", err)
	}

	if _, err := os.Stat(objectPath); err == nil {
		return nil
	}

	file, err := os.Create(objectPath)
	if err != nil {
		return fmt.Errorf("failed to create object file: %w", err)
	}
	defer file.Close()

	writer := zlib.NewWriter(file)
	defer writer.Close()

	if _, err := writer.Write(commit.Format()); err != nil {
		return fmt.Errorf("failed to write object content: %w", err)
	}

	return nil
}

func (db *Database) ReadCommit(hash string) (*objects.Commit, error) {
	if !utils.ValidateHash(hash) {
		return nil, fmt.Errorf("invalid hash: %s", hash)
	}

	objectPath := utils.ObjectPath(db.RepoPath, hash)

	file, err := os.Open(objectPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("object %s not found", hash)
		}
		return nil, fmt.Errorf("failed to open object file: %w", err)
	}
	defer file.Close()

	reader, err := zlib.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create zlib reader: %w", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read object content: %w", err)
	}

	nullIndex := -1
	for i, b := range content {
		if b == 0 {
			nullIndex = i
			break
		}
	}

	if nullIndex == -1 {
		return nil, fmt.Errorf("invalid commit format: missing null separator")
	}

	commitData := string(content[nullIndex+1:])

	lines := strings.Split(commitData, "\n")

	tree := ""
	parent := ""
	author := ""
	committer := ""
	message := ""

	messageStart := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "tree ") {
			tree = strings.TrimPrefix(line, "tree ")
		} else if strings.HasPrefix(line, "parent ") {
			parent = strings.TrimPrefix(line, "parent ")
		} else if strings.HasPrefix(line, "author ") {
			author = strings.TrimPrefix(line, "author ")
		} else if strings.HasPrefix(line, "committer ") {
			committer = strings.TrimPrefix(line, "committer ")
		} else if line == "" && messageStart == -1 {
			messageStart = i + 1
		}
	}

	if messageStart > 0 && messageStart < len(lines) {
		message = strings.Join(lines[messageStart:], "\n")
	}

	commit := &objects.Commit{
		Tree:      tree,
		Parent:    parent,
		Author:    author,
		Committer: committer,
		Message:   message,
		Hash:      hash,
	}

	return commit, nil
}