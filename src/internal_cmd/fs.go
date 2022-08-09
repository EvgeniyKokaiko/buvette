package internal_cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func ReadFile() string {
	path := GetWd()
	fullFileName := fmt.Sprintf("%s/%s", path, FILENAME)
	fmt.Println(fmt.Sprintf("[BUVETTE]: BuvetteFile PATH: %s", fullFileName))
	file, err := ioutil.ReadFile(fullFileName)
	if err != nil {
		fmt.Println("[BUVETTE]: Sorry! No BuvetteFile file here :(", path)
		return ""
	}
	return string(file)
}

func GetWd() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("[BUVETTE]: ERROR:", err)
	}
	return path
}

var selfStack = []string{}

func ManageLastProcesses() {
	pid := strconv.Itoa(os.Getpid())
	fPath := fmt.Sprintf("%s/.pid", GetWd())
	if !FileExists(fPath) {
		_, err := os.Create(fPath)
		if err != nil {
			fmt.Println("[BUVETTE]: ERROR: ", err)
		}
	}
	file, _ := ioutil.ReadFile(fPath)
	fileContent := string(file)
	if fileContent != "" && pid != strings.TrimSpace(fileContent) {
		fmt.Println(pid, fileContent)
		z, _ := strconv.Atoi(fileContent)
		process, _ := os.FindProcess(z)
		process.Kill()

	}
	ioutil.WriteFile(fPath, []byte(pid), 0777)
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
