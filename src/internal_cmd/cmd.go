package internal_cmd

import (
	"bufio"
	"buvette/src/doc"
	"buvette/src/types"
	"buvette/src/utils"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
	"syscall"
)

const FILENAME = "BuvetteFile"
const REGEXP = `(?m:(@[a-zA-Z0-9 ]*: (\[(.*?)\])*))(?:(: (\{(.*?)\}))*)`

var reservedFlags = []string{Help, Version, Author, Full, Current}

const (
	Help    = "--help"
	Version = "--version"
	Author  = "--author"
	Full    = "--full"
	Current = "--current"
)

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
	}
}

func (app *AppType) ManageStd(command []string, config map[string]string) {
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
	in := make(chan any)
	out := make(chan any)
	er := make(chan any)
	channels := map[string]any{
		"in":  in,
		"out": out,
		"er":  er,
	}
	scanner := bufio.NewScanner(stdout)
	errorScanner := bufio.NewScanner(stderr)
	stdinScanner := bufio.NewScanner(os.Stdin)
	//empty := " "
	go func(ch chan any) {
		for stdinScanner.Scan() {
			if utils.StandardizeSpaces(stdinScanner.Text()) == "r" {
				fmt.Println("[BUVETTE]: Reloading Process", fmt.Sprintf("`%s`", strings.Join(command, " ")))
				ch <- 54
			}
			if utils.StandardizeSpaces(stdinScanner.Text()) == "exit" {
				fmt.Println("[BUVETTE]: ExitApplication Process", fmt.Sprintf("`%s`", strings.Join(command, " ")))
				ch <- 55
			}
		}
	}(in)
	go func(ch chan any) {
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}(out)
	go func(ch chan any) {
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
				RestartSelf(config)
			} else if code == 55 {
				ExitApplication(config, channels)
			}
			break
		}

	}
}

func RestartSelf(config map[string]string) error {
	self, err := os.Executable()
	args := os.Args
	env := os.Environ()
	if err != nil {
		fmt.Println("[BUVETTE]: ERROR!4", err)
	}
	KillStdPort(config)
	go func() {
		cmd := exec.Command(self, args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = env
		err2 := cmd.Run()
		if err2 == nil {
			os.Exit(0)
		}
	}()
	return syscall.Exec(self, args, env)
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
