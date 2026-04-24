package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/AdeshDeshmukh/helix/internal/index"
	"github.com/AdeshDeshmukh/helix/internal/objects"
	"github.com/AdeshDeshmukh/helix/internal/refs"
	"github.com/AdeshDeshmukh/helix/internal/storage"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		message, _ := cmd.Flags().GetString("message")
		if message == "" {
			return fmt.Errorf("message required. Use: helix commit -m \"message\"")
		}

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		repoPath, err := findRepositoryRoot(wd)
		if err != nil {
			return err
		}

		db := storage.NewDatabase(repoPath)
		idx := index.NewIndex(repoPath)
		err = idx.Load()
		if err != nil {
			return err
		}

		if len(idx.GetEntries()) == 0 {
			return fmt.Errorf("nothing to commit")
		}

		fileCount := len(idx.GetEntries())

		treeBuilder := storage.NewTreeBuilder(db)
		for _, entry := range idx.GetEntries() {
			treeBuilder.AddEntry(entry.Path, entry.Hash, entry.Mode)
		}

		tree, err := treeBuilder.BuildTree()
		if err != nil {
			return err
		}

		err = db.WriteTree(tree)
		if err != nil {
			return err
		}

		var parentHash string
		parentHash, _ = refs.GetHEAD(repoPath)

		authorName := os.Getenv("GIT_AUTHOR_NAME")
		if authorName == "" {
			authorName = "Adesh Deshmukh"
		}

		authorEmail := os.Getenv("GIT_AUTHOR_EMAIL")
		if authorEmail == "" {
			authorEmail = "adeshkd123@gmail.com"
		}

		author := fmt.Sprintf("%s <%s> %d +0000", authorName, authorEmail, time.Now().Unix())

		commit := objects.NewCommit(tree.Hash, parentHash, author, author, message)

		err = db.WriteCommit(commit)
		if err != nil {
			return err
		}

		err = refs.SetHEAD(repoPath, commit.Hash)
		if err != nil {
			return err
		}

		mainBranch := "main"
		err = refs.SetBranchHash(repoPath, mainBranch, commit.Hash)
		if err != nil {
			return err
		}

		idx.Clear()
		err = idx.Save()
		if err != nil {
			return err
		}

		fmt.Printf("[main %s] %s\n", commit.Hash[:7], message)
		fmt.Printf("%d file(s) changed\n", fileCount)

		return nil
	},
}

func init() {
	commitCmd.Flags().StringP("message", "m", "", "commit message")
}
