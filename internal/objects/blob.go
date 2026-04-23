package objects

import (
	"fmt"
	"io"
	"os"

	"github.com/AdeshDeshmukh/helix/internal/utils"
)

type Blob struct {
	Content []byte
	Hash    string
	Size    int
}

func NewBlobFromFile(filePath string) (*Blob, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewBlobFromReader(file)
}

func NewBlobFromReader(reader io.Reader) (*Blob, error) {
	hash, content, err := utils.HashReader("blob", reader)
	if err != nil {
		return nil, err
	}
	return &Blob{
		Content: content,
		Hash:    hash,
		Size:    len(content),
	}, nil
}

func NewBlobFromContent(content []byte) *Blob {
	hash := utils.HashContent("blob", content)
	return &Blob{
		Content: content,
		Hash:    hash,
		Size:    len(content),
	}
}

func (b *Blob) Type() string {
	return "blob"
}

func (b *Blob) Format() []byte {
	header := fmt.Sprintf("blob %d\x00", len(b.Content))
	return append([]byte(header), b.Content...)
}

func (b *Blob) String() string {
	return fmt.Sprintf("Blob(%s, %d bytes)", b.Hash[:7], b.Size)
}
