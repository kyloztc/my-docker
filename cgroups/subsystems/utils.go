package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

// FindCgroupMountpoint 查找cgroup挂载点
func FindCgroupMountpoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}
	return ""
}

// GetCgroupPath 获取某一subsystem挂载点
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)
	fullCGroupPath := path.Join(cgroupRoot, cgroupPath)
	fmt.Printf("fill path: %v\n", fullCGroupPath)
	_, err := os.Stat(fullCGroupPath)
	if err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			err = os.Mkdir(fullCGroupPath, 0755)
			if err != nil {
				return "", fmt.Errorf("craete cgroup error|%v", err)
			}
		}
		return fullCGroupPath, nil
	}
	return "", fmt.Errorf("cgroup path error: %v", err)
}
