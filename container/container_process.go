package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func NewParentProcess(tty bool, command []string) (*exec.Cmd, *os.File) {
	fmt.Printf("new parent process|tty: %v|command: %v\n", tty, command)
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		fmt.Printf("new pipe error|%v\n", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	mntURL := "/root/mnt/"
	rootURL := "/root/"
	NewWorkSpace(rootURL, mntURL)
	cmd.Dir = mntURL
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

func NewWorkSpace(rootURL string, mntURL string) {
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)
}

func CreateReadOnlyLayer(rootURL string) {
	busyboxURL := fmt.Sprintf("%sbusybox/", rootURL)
	busyboxTarURL := fmt.Sprintf("%sbusybox.tar", rootURL)
	exist, err := PathExists(busyboxURL)
	if err != nil {
		fmt.Printf("check file path exist error|%v\n", err)
	}
	if !exist {
		if err := os.Mkdir(busyboxURL, 0777); err != nil {
			fmt.Printf("mkdir error|%v|path=%s\n", err, busyboxURL)
			return
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			fmt.Printf("untar dir %s error: %s", busyboxTarURL, err)
		}
	}
}

func CreateWriteLayer(rootURL string) {
	writeURL := fmt.Sprintf("%swriteLayer", rootURL)
	if err := os.Mkdir(writeURL, 0777); err != nil {
		fmt.Printf("mkdir error|%v|path: %v\n", err, writeURL)
	}
}

func CreateMountPoint(rootURL string, mntURL string) {
	if err := os.Mkdir(mntURL, 0777); err != nil {
		fmt.Printf("mkdir error|%v|path: %v", err, mntURL)
	}
	dirs := fmt.Sprintf("dirs=%swriteLayer:%sbusybox", rootURL, rootURL)
	fmt.Printf("dirs: %v\n", dirs)
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("create mount point error|%v\n", err)
	}
}

func DeleteWorkSpace(rootURL string, mntURL string) {
	DeleteMountPoint(rootURL, mntURL)
	DeleteWriteLayer(rootURL)
}

func DeleteMountPoint(rootURL string, mntURL string) {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("unmount error|%v\n", err)
	}
	if err := os.RemoveAll(mntURL); err != nil {
		fmt.Printf("remove error|%v\n", err)
	}
}

func DeleteWriteLayer(rootURL string) {
	writeURL := rootURL + "writeLayer/"
	if err := os.RemoveAll(writeURL); err != nil {
		fmt.Printf("Remove dir %s error %v", writeURL, err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, nil
}
