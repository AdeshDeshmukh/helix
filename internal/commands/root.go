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
	rootCmd.AddCommand(hashObjectCmd)
	rootCmd.AddCommand(writeTreeCmd)
	rootCmd.AddCommand(lsTreeCmd)
	rootCmd.AddCommand(commitTreeCmd)
	rootCmd.AddCommand(catFileCmd)
	rootCmd.AddCommand(logCmd)
}

func Execute(v, c, d string) error {
	version = v
	commit = c
	date = d
	return rootCmd.Execute()
}
