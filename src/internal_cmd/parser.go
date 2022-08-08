package internal_cmd

import (
	"buvette/src/types"
	"strings"
)

func (commands *AppCommands) ParseCommand(data []string) {
	path := GetWd()
	for _, node := range data {
		data := strings.Split(node, ":")
		entry := types.Command{}
		config := map[string]string{}
		if len(data) > 2 && len(data[2]) != 0 {
			parsedConfig := strings.Split(data[2][2:len(data[2])-1], ",")
			for _, subNode := range parsedConfig {
				if subNode != " " {
					parsedSubNode := strings.Split(subNode, "=")
					config[strings.TrimSpace(parsedSubNode[0])] = strings.TrimSpace(parsedSubNode[1])
				}
			}
		}
		config["__OriginalPath"] = path
		entry.Args = data[1]
		entry.Config = config
		COMMANDS[data[0]] = entry
	}
}
