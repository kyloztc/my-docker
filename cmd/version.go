package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("App demo version v1.0 -- HEAD")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
