package main

import "buvette/src/internal_cmd"

func main() {
	fileStr := internal_cmd.ReadFile()
	internal_cmd.RunApplication(fileStr)
}
