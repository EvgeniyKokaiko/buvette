package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

const FILENAME = "BuvetteFile"
const REGEXP = `(?m:(@[a-zA-Z0-9 ]*: (\[(.*?)\])*))`

const (
	Help    = "--help"
	Version = "--version"
	Author  = "--author"
	Full    = "--full"
	Current = "--current"
)

var reservedFlags = []string{Help, Version, Author, Full, Current}

func ReadFile() string {
	path, err := os.Getwd()
	if err != nil {
		panic("Error! on ReadFile ex")
	}
	fullFileName := fmt.Sprintf("%s/%s", path, FILENAME)
	fmt.Println(fullFileName)
	file, err := ioutil.ReadFile(fullFileName)
	if err != nil {
		fmt.Println("[BUVETTE]: Sorry! No BuvetteFile file here :(", path)
		return ""
	}
	return string(file)
}

var COMMANDS = map[string]string{}

func Runner(fileData string) {
	if fileData == "" {
		os.Exit(0)
		return
	}
	normalizedFileData := standardizeSpaces(fileData)
	re := regexp.MustCompile(REGEXP)
	index := re.FindAllString(normalizedFileData, 999)
	for _, node := range index {
		data := strings.Split(node, ":")
		COMMANDS[data[0]] = data[1]
	}
	argsWithProg := os.Args[len(os.Args)-1]
	isReserved := Contains[string](reservedFlags, strings.TrimSpace(argsWithProg))
	if isReserved {
		app.renderServiceInfo(argsWithProg)
		return
	}

	cmdAddon := COMMANDS[fmt.Sprintf("@%s", argsWithProg)]
	if cmdAddon == "" {
		fmt.Println("[BUVETTE]: Sorry! No such command :(")
		return
	}
	normalizeAddon := strings.TrimSpace(cmdAddon[2 : len(cmdAddon)-1])
	currentCommand := strings.Split(normalizeAddon, " ")

	fmt.Println("[BUVETTE]: CurrentCommand =", currentCommand)
	RunContext(currentCommand)
}

func (app *Application) renderServiceInfo(value string) {
	switch value {
	case Help:
		fmt.Println(app.HelpInfo)
		break
	case Version:
		fmt.Println(fmt.Sprintf("Buvette Version: %s", app.Version))
		break
	case Author:
		fmt.Println(fmt.Sprintf("Buvette Author: %s", app.Author))
		break
	case Full:
		fmt.Println(app)
		break
	case Current:
		for key, value := range COMMANDS {
			fmt.Println(fmt.Sprintf("Buvette Commands: %s", key), value)
		}
		break
	}
}

func RunContext(command []string) {
	cmd := exec.Command(strings.TrimSpace(command[0]), command[1:]...)
	stdout, errPipe := cmd.StdoutPipe()
	stderr, errStdErr := cmd.StderrPipe()
	stdin, _ := cmd.StdinPipe()
	inWriter := io.MultiWriter(stdin)
	go func() {
		io.Copy(inWriter, cmd.Stdin)
		defer stdin.Close()
	}()
	err := cmd.Start()
	if err != nil || errPipe != nil || errStdErr != nil {
		fmt.Println("[BUVETTE]: ERROR!", err)
	}
	in := make(chan int)
	out := make(chan string)
	er := make(chan string)
	scanner := bufio.NewScanner(stdout)
	errorScanner := bufio.NewScanner(stderr)
	stdinScanner := bufio.NewScanner(os.Stdin)
	//empty := " "
	go func(ch chan int) {
		for stdinScanner.Scan() {
			if standardizeSpaces(stdinScanner.Text()) == "r" {
				fmt.Println("[BUVETTE]: Reloading Process", strings.Join(command, ""))
				ch <- 54
			}
			if standardizeSpaces(stdinScanner.Text()) == "exit" {
				fmt.Println("[BUVETTE]: Exit Process", strings.Join(command, ""))
				ch <- 55
			}
		}
	}(in)
	go func(ch chan string) {
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}(out)
	go func(ch chan string) {
		for errorScanner.Scan() {
			ch <- errorScanner.Text()
		}
	}(er)
	for {
		select {
		case text := <-out:
			fmt.Println(text)
			break
		case errorMsg := <-er:
			fmt.Println(errorMsg)
			break
		case code := <-in:
			if code == 54 {
				RestartSelf("3000")
			} else if code == 55 {
				Exit("3000")
			}
			break
		}

	}
}

func RestartSelf(PORT string) error {
	self, err := os.Executable()
	args := os.Args
	env := os.Environ()
	if err != nil {
		fmt.Println("[BUVETTE]: ERROR!4", err)
	}
	KillStdPort(PORT)
	go func() {
		if runtime.GOOS == "windows" {
			cmd := exec.Command(self, args[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			cmd.Env = env
			err := cmd.Run()
			if err == nil {
				os.Exit(0)
			}
		}
	}()
	return syscall.Exec(self, args, env)
}

func KillStdPort(PORT string) {
	//powershell -ExecutionPolicy Bypass -Command "(Get-NetTCPConnection -LocalPort 3000).OwningProcess"
	checkPIDS := fmt.Sprintf("(Get-NetTCPConnection -LocalPort %s).OwningProcess", PORT)
	command := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", checkPIDS)
	stdout, errStdOut := command.StdoutPipe()
	err := command.Start()
	if err != nil {
		fmt.Println(err)
	}
	if errStdOut != nil {
		fmt.Println("[BUVETTE]: ERROR!1", errStdOut)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		kill := exec.Command("taskkill", "/T", "/F", "/PID", strings.TrimSpace(scanner.Text()))
		err123 := kill.Run()
		if err123 != nil {
			fmt.Println("[BUVETTE]: ERROR!3", err123)
		}
	}
}

func Exit(PORT string) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		KillStdPort(PORT)
		wg.Done()
	}()
	os.Exit(0)
}
