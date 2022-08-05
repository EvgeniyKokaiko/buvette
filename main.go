package main

type Application struct {
	Name     string
	Usage    string
	Author   string
	Version  string
	HelpInfo string
}

var app = Application{
	Name:    "Buvette",
	Usage:   "MakeFile analog",
	Author:  "re1nhart",
	Version: "0.0.9",
	HelpInfo: ` 
				Help    = "--help"
				Version = "--version"
				Author  = "--author"
				Full    = "--full"
				Current = "--current"
				To reload, write r several times and wait.`,
}

func main() {
	fileStr := ReadFile()
	Runner(fileStr)
}
