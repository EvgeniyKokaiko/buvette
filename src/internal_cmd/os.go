package internal_cmd

import (
	"bufio"
	"buvette/src/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func ExecCommandForPort(PORT string) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		//(Get-Process -Name src).Id
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

func RunCallbacksByPlatform(winCallback func(), macosCallback func()) {
	switch runtime.GOOS {
	case "windows":
		winCallback()
		break
	case "darwin":
		macosCallback()
		break
	default:
		winCallback()
		break
	}
}

func KillProcessByHandles() {
	processName := strings.Split(filepath.Base(os.Args[0]), ".")[0]
	//(Get-Process -Name src) | Select Id, Handles
	//(Get-Process -Name src) | Select Id, Handles | Sort -Property Handles
	pids := []string{}
	checkPIDS := fmt.Sprintf("(Get-Process -Name %s) | Select Id, Handles | Sort -Property Handles", processName)
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", checkPIDS)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	reader := bufio.NewScanner(stdout)
	for reader.Scan() {
		pids = append(pids, utils.StandardizeSpaces(strings.TrimSpace(reader.Text())))
	}
	services := pids[3:]
	src := []string{}
	for _, v := range services {
		source := strings.Split(v, " ")
		for _, p := range source {
			if v != "" {
				src = append(src, p)
			}
		}
	}
	lastPid := src[len(src)-2]
	for k, v := range src {
		if len(src) > 1 && (k%2 != 0 || v == lastPid) {
			continue
		} else {
			fmt.Println("[BUVETTE]: Process with PID:", v, "was eliminated!")
			kill := exec.Command("taskkill", "/F", "/PID", strings.TrimSpace(v))
			err123 := kill.Start()
			if err123 != nil {
				fmt.Println("[BUVETTE]: ERROR!3", err123)
			}
		}
	}
}
