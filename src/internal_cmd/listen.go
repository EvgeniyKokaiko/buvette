package internal_cmd

import (
	"bufio"
	"buvette/src/utils"
	"fmt"
	"io"
	"os"
	"strings"
)

func ListenProcesses(command []string, config map[string]string, stdin io.Reader, stdout io.Writer, stdinErr io.Writer, stdinR io.WriteCloser, stdoutR io.ReadCloser, stderrR io.ReadCloser) {
	inWriter := io.MultiWriter(stdinR)
	go func() {
		io.Copy(inWriter, stdin)
		defer stdinR.Close()
	}()

	in := make(chan any)
	out := make(chan any)
	er := make(chan any)
	channels := map[string]any{
		"in":  in,
		"out": out,
		"er":  er,
	}
	scanner := bufio.NewScanner(stdoutR)
	errorScanner := bufio.NewScanner(stderrR)
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
				StartDaemon(command, config)
			} else if code == 55 {
				ExitApplication(config, channels)
			}
			break
		default:
		}
	}
}
