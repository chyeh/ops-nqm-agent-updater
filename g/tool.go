package g

import (
	"fmt"
	"os/exec"
	"strings"
)

func CheckModulePid(name string) (string, error) {
	output, err := exec.Command("pgrep", "-f", name).Output()
	if err != nil {
		return "", err
	}
	pidStr := strings.TrimSpace(string(output))
	return pidStr, nil
}

func CheckModuleStatus(name string) int {
	fmt.Print("Checking status [", name, "]...")

	pidStr, err := CheckModulePid(name)
	if err != nil {
		fmt.Println("not running!!")
		return NotRunning
	}

	fmt.Println("running with PID [", pidStr, "]!!")
	return Running
}
