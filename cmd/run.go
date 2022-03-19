package cmd

import (
	"fmt"
	"os"
	"strings"

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
		Run(true, args)
	},
}

func Run(tty bool, command []string) {
	fmt.Printf("Run|tty: %v|command: %v", tty, command)
	parent, writePipe := container.NewParentProcess(tty, command)
	if parent == nil {
		fmt.Printf("new parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		fmt.Printf("run error|%v", err)
	}
	sendInitCommand(command, writePipe)
	parent.Wait()
	fmt.Printf("my docker exit")
	os.Exit(-1)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	fmt.Printf("command all is %s\n", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func init() {
	rootCmd.AddCommand(runCmd)
}
