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
	Version: "0.1.0",
	HelpInfo: ` 
				Help    = "--help"
				Version = "--version"
				Author  = "--author"
				Full    = "--full"
				Current = "--current"
				To reload, write r and press enter. To exit you can write exit and press enter`,
}

func main() {
	fileStr := ReadFile()
	Runner(fileStr)
}
