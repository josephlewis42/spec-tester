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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tidyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tidyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
