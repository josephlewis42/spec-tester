package cmd

import (
	"fmt"

	"github.com/josephlewis42/scheme-compliance/tester/model/storage"
	"github.com/josephlewis42/scheme-compliance/tester/validation"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check path/to/spec",
	Short: "Check the specification for errors.",
	Long:  `Performs tests for the specification at the given path.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		cmd.SilenceUsage = true

		suite, err := storage.LoadSuite(args[0])
		if err != nil {
			return err
		}

		fmt.Fprintf(
			cmd.OutOrStdout(),
			"Loaded %d specs, %d implementations, %d tests\n",
			len(suite.Specifications),
			len(suite.Implementations),
			len(suite.Tests),
		)

		var vs validation.ValidationSummary
		suite.RunValidation(func(name string, v *validation.Validator) {
			fmt.Fprintf(cmd.OutOrStdout(), "Checking: %s\n", name)

			for _, result := range v.Results {
				fmt.Fprintln(cmd.OutOrStdout(), "-", result.String())
			}

			vs.Update(v)
		})

		fmt.Fprintf(
			cmd.OutOrStdout(),
			"\nResults: %d Errors, %d Warnings, %d Infos\n",
			vs.ErrorCount,
			vs.WarningCount,
			vs.InfoCount,
		)

		if vs.ErrorCount > 0 {
			return fmt.Errorf("%d errors encountered", vs.ErrorCount)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
