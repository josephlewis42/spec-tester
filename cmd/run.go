package cmd

import (
	"github.com/josephlewis42/scheme-compliance/tester/executor"
	"github.com/josephlewis42/scheme-compliance/tester/model/storage"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tests",
	Long:  `Run tests for all specifications.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		suite, err := storage.LoadSuite(args[0])
		if err != nil {
			return err
		}

		return executor.Execute(cmd.Context(), suite, executor.ExecutionOptions{})
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
