package cmd

import (
	"fmt"

	"github.com/just-hms/dow/internal"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "displays the version of dow",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(internal.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
