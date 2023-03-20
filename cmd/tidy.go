package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// tidyCmd represents the tidy command
var tidyCmd = &cobra.Command{
	Use:   "tidy",
	Short: "Organizes the tests under the given directory.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tidy called")
	},
}

func init() {
	rootCmd.AddCommand(tidyCmd)
}
