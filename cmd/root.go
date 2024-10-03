package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "dow",
	Short:   "move the last downloaded file in the current (or the specified) folder",
	Example: "use `dow` or `dow /target/folder` to move the latest file where you want",
	RunE: func(cmd *cobra.Command, args []string) error {
		return mvCmd.RunE(cmd, args)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true

	rootCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "show the name of the moved file")
	rootCmd.Flags().BoolVarP(&yesFlag, "yes", "y", false, "force dow to move the latest file even if it's old")
}
