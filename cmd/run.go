package cmd

import (
	"fmt"
	"os"

	"my-docker/container"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("run cmd args: %v\n", args)
		if len(args) < 1 {
			fmt.Printf("missing command: %v\n", args)
		}
		command := args[0]
		Run(true, command)
	},
}

func Run(tty bool, command string) {
	fmt.Printf("Run|tty: %v|command: %v", tty, command)
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		fmt.Printf("run error|%v", err)
	}
	parent.Wait()
	fmt.Printf("my docker exit")
	os.Exit(-1)
}

func init() {
	rootCmd.AddCommand(runCmd)
}
