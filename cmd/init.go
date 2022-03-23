package cmd

import (
	"fmt"

	"my-docker/container"

	"github.com/spf13/cobra"
)

// 初始化容器命令
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init cmd",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("init args: %v\n", args)
		err := container.RunContainerInitProcess()
		if err != nil {
			fmt.Printf("init cmd run error|%v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
