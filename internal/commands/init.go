package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a new Helix repository",
	Long: `Create a new .helix repository in the current or specified directory.

Structure created:
  .helix/
  ├── objects/      # Object database
  ├── refs/
  │   ├── heads/    # Branches
  │   └── tags/     # Tags
  ├── HEAD          # Current branch
  └── config        # Configuration
`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := initializeRepository(absPath); err != nil {
		return err
	}

	fmt.Printf("✨ Initialized empty Helix repository in %s/.helix/\n", absPath)
	return nil
}

func initializeRepository(path string) error {
	helixDir := filepath.Join(path, ".helix")

	if _, err := os.Stat(helixDir); err == nil {
		return fmt.Errorf("repository already exists at %s", helixDir)
	}

	dirs := []string{
		helixDir,
		filepath.Join(helixDir, "objects"),
		filepath.Join(helixDir, "refs"),
		filepath.Join(helixDir, "refs", "heads"),
		filepath.Join(helixDir, "refs", "tags"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", dir, err)
		}
	}

	headPath := filepath.Join(helixDir, "HEAD")
	if err := os.WriteFile(headPath, []byte("ref: refs/heads/main\n"), 0644); err != nil {
		return fmt.Errorf("failed to create HEAD: %w", err)
	}

	configPath := filepath.Join(helixDir, "config")
	configContent := "[core]\n\trepositoryformatversion = 0\n\tfilemode = true\n\tbare = false\n"
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	return nil
}
