package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version string
var revision string

// version command shows version and build revision
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version and build revision.",
	Long:  `Show version and build revision.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s-%s\n", version, revision)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
