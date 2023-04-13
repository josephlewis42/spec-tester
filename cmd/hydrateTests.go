package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/josephlewis42/scheme-compliance/tester/model/storage"
	"github.com/spf13/cobra"
)

// hydrateTestsCmd represents the hydrateTests command
var hydrateTestsCmd = &cobra.Command{
	Use:   "hydrate-tests path/to/spec",
	Short: "Display the hydreated tests for the given path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		suite, err := storage.LoadSuite(args[0])
		if err != nil {
			return err
		}

		hydrated := suite.ListTests()
		out, err := json.MarshalIndent(hydrated, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(out))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(hydrateTestsCmd)
}
