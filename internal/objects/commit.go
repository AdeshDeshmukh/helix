package objects

import (
	"fmt"
	"strings"

	"github.com/AdeshDeshmukh/helix/internal/utils"
)

type Commit struct {
	Tree      string
	Parent    string
	Author    string
	Committer string
	Message   string
	Hash      string
}

func NewCommit(tree, parent, author, committer, message string) *Commit {
	commit := &Commit{
		Tree:      tree,
		Parent:    parent,
		Author:    author,
		Committer: committer,
		Message:   message,
	}
	commit.Hash = utils.HashContent("commit", commit.Format())
	return commit
}

func (c *Commit) Type() string {
	return "commit"
}

func (c *Commit) Format() []byte {
	var content strings.Builder

	content.WriteString(fmt.Sprintf("tree %s\n", c.Tree))

	if c.Parent != "" {
		content.WriteString(fmt.Sprintf("parent %s\n", c.Parent))
	}

	content.WriteString(fmt.Sprintf("author %s\n", c.Author))
	content.WriteString(fmt.Sprintf("committer %s\n", c.Committer))
	content.WriteString("\n")
	content.WriteString(c.Message)

	contentBytes := []byte(content.String())
	header := fmt.Sprintf("commit %d\x00", len(contentBytes))

	return append([]byte(header), contentBytes...)
}

func (c *Commit) String() string {
	return fmt.Sprintf("Commit(%s)", c.Hash[:7])
}
