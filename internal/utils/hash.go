
package utils

import (
	"crypto/sha1"
	"fmt"
	"io"
)

func HashContent(objectType string, content []byte) string {
	header := fmt.Sprintf("%s %d\x00", objectType, len(content))
	h := sha1.New()
	h.Write([]byte(header))
	h.Write(content)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HashReader(objectType string, reader io.Reader) (string, []byte, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", nil, err
	}
	hash := HashContent(objectType, content)
	return hash, content, nil
}

func ObjectPath(repoPath, hash string) string {
	return fmt.Sprintf("%s/.helix/objects/%s/%s", repoPath, hash[:2], hash[2:])
}

func ValidateHash(hash string) bool {
	if len(hash) != 40 {
		return false
	}
	for _, c := range hash {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

func HexToBytes(hex string) ([]byte, error) {
	if len(hex) != 40 {
		return nil, fmt.Errorf("invalid hex length: expected 40, got %d", len(hex))
	}

	bytes := make([]byte, 20)
	for i := 0; i < 40; i += 2 {
		var b byte
		for j := 0; j < 2; j++ {
			c := hex[i+j]
			var v byte
			if c >= '0' && c <= '9' {
				v = c - '0'
			} else if c >= 'a' && c <= 'f' {
				v = c - 'a' + 10
			} else if c >= 'A' && c <= 'F' {
				v = c - 'A' + 10
			} else {
				return nil, fmt.Errorf("invalid hex character: %c", c)
			}
			b = b*16 + v
		}
		bytes[i/2] = b
	}
	return bytes, nil
}

func BytesToHex(bytes []byte) string {
	hex := make([]byte, len(bytes)*2)
	for i, b := range bytes {
		hex[i*2] = "0123456789abcdef"[b>>4]
		hex[i*2+1] = "0123456789abcdef"[b&0x0f]
	}
	return string(hex)
}

func GetFileMode(isDir bool, isExecutable bool) string {
	if isDir {
		return "040000"
	}
	if isExecutable {
		return "100755"
	}
	return "100644"
}