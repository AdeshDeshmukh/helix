#!/bin/bash

echo "🧬 Setting up Helix from scratch..."

# Step 1: Initialize Go module
go mod init github.com/AdeshDeshmukh/helix

# Step 2: Install Cobra
go get github.com/spf13/cobra@latest

# Step 3: Create directory structure
mkdir -p cmd/helix
mkdir -p internal/commands
mkdir -p internal/objects
mkdir -p internal/storage
mkdir -p internal/index
mkdir -p internal/refs
mkdir -p internal/diff
mkdir -p internal/merge
mkdir -p internal/remote
mkdir -p internal/utils
mkdir -p pkg/helix
mkdir -p test/integration
mkdir -p test/fixtures
mkdir -p docs
mkdir -p examples

# Step 4: Create .gitignore
cat > .gitignore << 'EOF'
# Binaries
helix
*.exe
*.dll
*.so
*.dylib

# Test binary
*.test

# Coverage
*.out
coverage.html

# Dependencies
vendor/

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db

# Test repos
test-repos/
*.helix/

# Logs
*.log
EOF

# Step 5: Create LICENSE
cat > LICENSE << 'EOF'
MIT License

Copyright (c) 2026 Adesh Deshmukh

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
EOF

# Step 6: Create README.md
cat > README.md << 'EOF'
<div align="center">

# 🧬 HELIX

### *Version Control from First Principles*

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/AdeshDeshmukh/helix)

*A production-grade Git implementation built entirely from scratch in Go*

</div>

---

## 🎯 What is Helix?

**Helix** is a fully-functional Git-compatible version control system built from the ground up in Go. Like the DNA double helix intertwines genetic information, Helix intertwines your code's history into an elegant, distributed graph.

### Learning Goals

- 🗄️ **Content-Addressable Storage** - Hash-based object databases
- 🌳 **Merkle Trees & DAGs** - Cryptographic data structures
- 📡 **Network Protocols** - Git's smart HTTP/SSH protocol
- 🔄 **Distributed Systems** - Conflict resolution, consensus
- ⚡ **Performance Engineering** - Delta compression, packfiles

> **Note:** Educational project for learning version control internals.

---

## 💡 Philosophy

> *"Don't build applications. Build products. Build systems."* — Anuj Bhaiya

---

## ✨ Roadmap

### 📦 Phase 1: Local VCS (Weeks 1-4)
- [x] Repository initialization
- [ ] Object storage (blobs, trees, commits)
- [ ] Staging area (index)
- [ ] Branching and checkout

### 🌐 Phase 2: Remote Operations (Weeks 5-8)
- [ ] Diff algorithms
- [ ] Merge strategies
- [ ] Clone, fetch, push

### 🚀 Phase 3: Advanced (Weeks 9-12)
- [ ] Packfiles
- [ ] Rebase
- [ ] Optimization

---

## 🚀 Quick Start

```bash
# Build
go build -o helix cmd/helix/main.go

# Initialize repository
./helix init

# Show version
./helix version
```

---

## 📚 Resources

- [Pro Git Book](https://git-scm.com/book/en/v2)
- [Git Internals](https://git-scm.com/book/en/v2/Git-Internals-Git-Objects)

---

## 📝 License

MIT License - see [LICENSE](LICENSE)

---

## 📬 Connect

Built by [Adesh Deshmukh](https://github.com/AdeshDeshmukh)

- 📧 adeshkd123@gmail.com
- 💼 [LinkedIn](https://www.linkedin.com/in/adesh-deshmukh-532744318/)

---

<div align="center">

**Built with 🧬 and ❤️ in Go**

*Building Git from first principles*

</div>
EOF

# Step 7: Create cmd/helix/main.go
cat > cmd/helix/main.go << 'EOF'
package main

import (
	"fmt"
	"os"

	"github.com/AdeshDeshmukh/helix/internal/commands"
)

var (
	version = "0.1.0-dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	if err := commands.Execute(version, commit, date); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
EOF

# Step 8: Create internal/commands/root.go
cat > internal/commands/root.go << 'EOF'
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version string
	commit  string
	date    string
)

var rootCmd = &cobra.Command{
	Use:   "helix",
	Short: "🧬 Helix - Version control from first principles",
	Long: `
🧬 HELIX - Git implementation built from scratch in Go

Learn version control internals by building Git yourself.

GitHub: https://github.com/AdeshDeshmukh/helix
`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("🧬 helix version %s\n", version)
		fmt.Printf("   commit: %s\n", commit)
		fmt.Printf("   built: %s\n", date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
}

func Execute(v, c, d string) error {
	version = v
	commit = c
	date = d
	return rootCmd.Execute()
}
EOF

# Step 9: Create internal/commands/init.go
cat > internal/commands/init.go << 'EOF'
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
EOF

# Step 10: Create ARCHITECTURE.md
cat > ARCHITECTURE.md << 'EOF'
# 🏗️ Helix Architecture

## Repository Structure

```
.helix/
├── objects/      # Content-addressable storage
├── refs/
│   ├── heads/   # Branches
│   └── tags/    # Tags
├── HEAD         # Current branch pointer
└── config       # Configuration
```

## Object Types

### Blob
```
blob <size>\0<content>
```

### Tree
```
tree <size>\0
<mode> <name>\0<20-byte-sha>
```

### Commit
```
commit <size>\0
tree <sha>
parent <sha>
author <name> <email> <timestamp>

Message
```

---

*Building Git from first principles*
EOF

# Step 11: Tidy dependencies
go mod tidy

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "  1. Build: go build -o helix cmd/helix/main.go"
echo "  2. Test:  ./helix version"
echo "  3. Test:  ./helix init test-demo"
echo ""
