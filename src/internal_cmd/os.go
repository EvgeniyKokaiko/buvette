package internal_cmd

import (
	"fmt"
	"os/exec"
	"runtime"
)

func ExecCommandForPort(PORT string) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		//powershell -ExecutionPolicy Bypass -Command "(Get-NetTCPConnection -LocalPort 3000).OwningProcess"
		checkPIDS := fmt.Sprintf("(Get-NetTCPConnection -LocalPort %s).OwningProcess", PORT)
		return exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", checkPIDS)
	case "darwin":
		pid := fmt.Sprintf("tcp:%s", PORT)
		return exec.Command("lsof", "-t", "-i", pid, "|", "xargs", "kill")
	default:
		checkPIDS := fmt.Sprintf("(Get-NetTCPConnection -LocalPort %s).OwningProcess", PORT)
		return exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", checkPIDS)
	}
}
