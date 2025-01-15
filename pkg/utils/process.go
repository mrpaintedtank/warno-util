package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func StopProcess(name string) error {
	cmd := exec.Command("taskkill", "/F", "/IM", name+".exe")
	err := cmd.Run()
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}
		return fmt.Errorf("an error occurred while stopping the '%s' process: %v", name, err)
	}
	return nil
}

func ProcessRunning(name string) bool {
	output, err := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s.exe", name)).Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), name+".exe")
}
