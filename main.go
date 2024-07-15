package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Command struct {
	Name   string
	Args   []string
	Flags  map[string]string
	Sub    []*Command
	Parent *Command
}

func parseYaml(yamlConfig map[string]interface{}) Command {
	var rootStr string
	if len(yamlConfig) != 1 {
		panic("Root command must be only one")
	}
	for key := range yamlConfig {
		rootStr = key
	}
	root := Command{
		Name:  rootStr,
		Flags: make(map[string]string),
	}

	for key, value := range yamlConfig {
		parseNode(key, value, &root)
	}

	return root
}

func parseNode(name string, value interface{}, parent *Command) {
	switch v := value.(type) {
	case []interface{}:
		for _, item := range v {
			switch item := item.(type) {
			case string:
				if item[0] == '+' {
					parent.Args = append(parent.Args, item[1:])
				} else if item[0] == '-' {
					parent.Flags[item] = ""
				} else {
					sub := Command{
						Name:   item,
						Flags:  make(map[string]string),
						Parent: parent,
					}
					parseNode(item, nil, &sub)
					parent.Sub = append(parent.Sub, &sub)
				}
			case map[interface{}]interface{}:
				for subCommand, subValue := range item {
					sub := Command{
						Name:   subCommand.(string),
						Flags:  make(map[string]string),
						Parent: parent,
					}
					parseNode(subCommand.(string), subValue, &sub)
					parent.Sub = append(parent.Sub, &sub)
				}
			}
		}
	case map[interface{}]interface{}:
		for subCommand, subValue := range v {
			sub := Command{
				Name:   subCommand.(string),
				Flags:  make(map[string]string),
				Parent: parent,
			}
			parseNode(subCommand.(string), subValue, &sub)
			parent.Sub = append(parent.Sub, &sub)
		}
	}
}

func main() {
	var yamlConfig map[string]interface{}
	f := "config.yaml"
	yamlData, err := os.ReadFile(f)
	err = yaml.Unmarshal(yamlData, &yamlConfig)
	if err != nil {
		panic(err)
	}

	commands := parseYaml(yamlConfig)
	printCommands(commands, 0)
}

func printCommands(cmd Command, level int) {
	for i := 0; i < level; i++ {
		print("  ")
	}
	println("Command:" + namePrint(cmd))
	for i := 0; i < level; i++ {
		print("  ")
	}
	println("Args:")
	for _, arg := range cmd.Args {
		for i := 0; i < level; i++ {
			print("  ")
		}
		println("  ", arg)
	}
	for i := 0; i < level; i++ {
		print("  ")
	}
	println("Flags:")
	for flag, value := range cmd.Flags {
		for i := 0; i < level; i++ {
			print("  ")
		}
		println("  ", flag, value)
	}
	for _, sub := range cmd.Sub {
		printCommands(*sub, level+1)
	}
}

func namePrint(cmd Command) string {
	// recursive print
	if cmd.Parent != nil {
		return namePrint(*cmd.Parent) + " " + cmd.Name
	}
	return cmd.Name
}

func MainTemplate() []byte {
	return []byte(`
package main

import "{{ .PkgName }}/cmd"

func main() {
	cmd.Execute()
}
`)
}
