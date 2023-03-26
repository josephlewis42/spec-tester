package cmd

import (
	"github.com/josephlewis42/scheme-compliance/tester/model/storage"
	"github.com/spf13/cobra"
)

// tidyCmd represents the tidy command
var tidyCmd = &cobra.Command{
	Use:   "tidy path/to/spec",
	Short: "Organizes the tests under the given directory.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		suite, err := storage.LoadSuite(args[0])
		if err != nil {
			return err
		}

		suite.Tidy()

		suite.Diff(cmd.OutOrStdout())
		suite.Save()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tidyCmd)
}
