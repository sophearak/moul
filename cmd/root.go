package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "moul",
	Short: "The minimalist photos collection generator",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("hey")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}