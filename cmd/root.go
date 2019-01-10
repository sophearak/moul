package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [string folder name]",
	Short: "Create new photo collection",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(args[0]); !os.IsNotExist(err) {
			fmt.Println("Folder already exists")
		}

		err := os.MkdirAll(args[0], os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		folders := []string{args[0] + "/photos", args[0] + "/photos/cover", args[0] + "/photos/profile", args[0] + "/photos/collection"}

		for _, folder := range folders {
			os.MkdirAll(folder, os.ModePerm)
		}

	},
}

func Execute() {
	var rootCmd = &cobra.Command{Use: "moul", Short: "The minimalist photo collection generator"}
	rootCmd.AddCommand(newCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
