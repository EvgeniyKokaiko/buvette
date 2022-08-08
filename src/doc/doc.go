package doc

import "buvette/src/types"

var App = types.Application{
	Name:    "Buvette",
	Usage:   "MakeFile analog",
	Author:  "re1nhart",
	Version: "0.1.4",
	HelpInfo: ` 
				Help    = "--help"
				Version = "--version"
				Author  = "--author"
				Full    = "--full"
				Current = "--current"
				To reload, write r and press enter. To exit you can write exit and press enter`,
}
