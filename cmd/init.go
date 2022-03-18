package cmd

import (
	"fmt"

	"my-docker/container"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init cmd",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("init args: %v\n", args)
		command := args[0]
		fmt.Printf("commond: %v\n", command)
		err := container.RunContainerInitProcess(command, args)
		if err != nil {
			fmt.Printf("init cmd run error|%v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
