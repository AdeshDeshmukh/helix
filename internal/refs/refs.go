package refs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetHEAD(repoPath string) (string, error) {
	headPath := filepath.Join(repoPath, ".helix", "HEAD")
	
	data, err := os.ReadFile(headPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	
	content := strings.TrimSpace(string(data))
	
	if strings.HasPrefix(content, "ref: ") {
		branch := strings.TrimPrefix(content, "ref: ")
		return GetBranchHash(repoPath, branch)
	}
	
	return content, nil
}

func SetHEAD(repoPath string, hash string) error {
	headPath := filepath.Join(repoPath, ".helix", "HEAD")
	
	err := os.WriteFile(headPath, []byte(hash+"\n"), 0644)
	if err != nil {
		return fmt.Errorf("failed to write HEAD: %w", err)
	}
	
	return nil
}

func GetCurrentBranch(repoPath string) (string, error) {
	headPath := filepath.Join(repoPath, ".helix", "HEAD")
	
	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}
	
	content := strings.TrimSpace(string(data))
	
	if strings.HasPrefix(content, "ref: ") {
		branch := strings.TrimPrefix(content, "ref: ")
		branch = strings.TrimPrefix(branch, "refs/heads/")
		return branch, nil
	}
	
	return "detached", nil
}

func GetBranchHash(repoPath string, branch string) (string, error) {
	if !strings.HasPrefix(branch, "refs/") {
		branch = "refs/heads/" + branch
	}
	
	branchPath := filepath.Join(repoPath, ".helix", branch)
	
	data, err := os.ReadFile(branchPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	
	return strings.TrimSpace(string(data)), nil
}

func SetBranchHash(repoPath string, branch string, hash string) error {
	if !strings.HasPrefix(branch, "refs/") {
		branch = "refs/heads/" + branch
	}
	
	branchPath := filepath.Join(repoPath, ".helix", branch)
	branchDir := filepath.Dir(branchPath)
	
	err := os.MkdirAll(branchDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create branch directory: %w", err)
	}
	
	err = os.WriteFile(branchPath, []byte(hash+"\n"), 0644)
	if err != nil {
		return fmt.Errorf("failed to write branch: %w", err)
	}
	
	return nil
}
