package container

import (
	"fmt"
	"os"
	"syscall"
)

func RunContainerInitProcess(command string, args []string) error {
	fmt.Printf("run container command: %v\n", command)
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		fmt.Printf("syscall mount error|%v\n", err)
	}
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		fmt.Printf("init process error|%v", err)
		return err
	}
	fmt.Printf("init process finish.")
	return nil
}
