package internal_cmd

import (
	"bufio"
	"buvette/src/doc"
	"buvette/src/types"
	"buvette/src/utils"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

const FILENAME = "BuvetteFile"
const REGEXP = `(?m:(@[a-zA-Z0-9 ]*: (\[(.*?)\])*))(?:(: (\{(.*?)\}))*)`

var reservedFlags = []string{Help, Version, Author, Full, Current, Example}

const (
	Help    = "--help"
	Version = "--version"
	Author  = "--author"
	Full    = "--full"
	Current = "--current"
	Example = "--example"
)

//Система такова!
//при старті аппки, ми зразу виконуємо os.StartProcess
// коли нажали r, тоді убиваємо його нахуй і по новій

type AppType struct {
	types.Application
}

type AppCommands map[string]types.Command

var App = AppType{doc.App}

var COMMANDS = AppCommands{}

func RunApplication(fileData string) {
	if fileData == "" {
		os.Exit(0)
		return
	}
	normalizedFileData := utils.StandardizeSpaces(fileData)
	re := regexp.MustCompile(REGEXP)
	index := re.FindAllString(normalizedFileData, 999)
	COMMANDS.ParseCommand(index)
	argsWith := os.Args[len(os.Args)-1]
	isReserved := utils.Contains[string](reservedFlags, strings.TrimSpace(argsWith))
	if isReserved {
		App.RenderServiceInfo(argsWith)
		return
	}

	cmdAddon := COMMANDS[fmt.Sprintf("@%s", argsWith)]
	if reflect.ValueOf(cmdAddon).IsZero() {
		fmt.Println("[BUVETTE]: Sorry! No such command :(")
		return
	}
	normalizeAddon := strings.TrimSpace(cmdAddon.Args[2 : len(cmdAddon.Args)-1])
	currentCommand := strings.Split(normalizeAddon, " ")
	fmt.Println("[BUVETTE]: CurrentCommand =", currentCommand)
	App.ManageStd(currentCommand, cmdAddon.Config)
}

func (app *AppType) RenderServiceInfo(value string) {
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
	case Example:
		fmt.Println(doc.Example())
		break
	default:
		fmt.Println("[BUVETTE]: Sorry, no such reserved flag!")
		break
	}
}

func (app *AppType) ManageStd(command []string, config map[string]string) {
	cdPath := config["PATH"]
	if cdPath != "" {
		err := os.Chdir(fmt.Sprintf("%s%s", GetWd(), cdPath[1:len(cdPath)-1]))
		if err != nil {
			fmt.Println("[BUVETTE]: ERROR!", err)
		}
	}
	StartDaemon(command, config)
}

func StartDaemon(command []string, config map[string]string) {
	KillStdPort(config)
	if runtime.GOOS == "windows" {
		env := os.Environ()
		cmd := exec.Command(strings.TrimSpace(command[0]), command[1:]...)
		stdout, errPipe := cmd.StdoutPipe()
		stderr, errStdErr := cmd.StderrPipe()
		stdin, _ := cmd.StdinPipe()
		if errPipe != nil || errStdErr != nil {
			fmt.Println("[BUVETTE]: ERROR!", errPipe)
		}
		cmd.Env = env
		//KillProcessByHandles()
		err := cmd.Start()
		ListenProcesses(command, config, cmd.Stdin, cmd.Stdout, cmd.Stderr, stdin, stdout, stderr)
		if err == nil {
			os.Exit(0)
		}
	}
	//doesn't work on Windows!
	//return syscall.Exec(self, args, env)
}

func KillStdPort(config map[string]string) {
	PORT := config["PORT"]
	if PORT == "" {
		fmt.Println(fmt.Sprintf("[BUVETTE]: Port of your application can not be killed because it's not specified in BuvetteFile!"))
		return
	}
	command := ExecCommandForPort(PORT)
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
	fmt.Println(fmt.Sprintf("[BUVETTE]: PORT %s is erased!", PORT))
}

func ExitApplication(config map[string]string, chans map[string]any) {
	isExit := make(chan bool)
	go func(exit chan bool) {
		KillStdPort(config)
		isExit <- true
		defer close(exit)
	}(isExit)
	for _, v := range chans {
		close(v.(chan any))
	}
	data := <-isExit
	if data {
		os.Exit(0)
	}
}
