package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v2"
)

type CLIConfig struct {
	Cmd map[string]interface{} `yaml:"cmd"`
}

const configFile = "config.yaml"

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("No commands provided.")
		return
	}

	config := loadConfig()

	keys := []string{}
	values := []string{}
	for _, arg := range args {
		if arg[0] == '-' || arg[0] == '+' {
			values = append(values, arg)
		} else {
			keys = append(keys, arg)
		}
	}

	//current := buildCommand(keys, values)

	//	currentm1 := buildMap(current)

	//	y1, _ := yaml.Marshal(currentm1)
	//fmt.Println("New command:")
	//fmt.Printf("%s\n", y1)

	cmd := parseYaml(config.Cmd)

	buildCommand2(&cmd, keys, values)
	resm := buildMap(cmd)

	config.Cmd = resm

	saveConfig(config)
}
func buildCommand2(root *Command, keys, values []string) {
	current := root
	depth := 0
	found := false

	for _, c := range current.Sub {
		if len(keys) <= depth {
			break
		}
		if c.Name == keys[depth] {
			found = true
			if len(keys)-1 <= depth {
				for _, value := range values {
					if value[0] == '-' {
						if _, ok := c.Flags[strings.TrimPrefix(strings.TrimPrefix(value, "-"), "-")]; !ok {
							c.Flags[strings.TrimPrefix(strings.TrimPrefix(value, "-"), "-")] = ""
						}
					} else {
						c.Args = append(c.Args, value[1:])
					}
				}
				return
			}
			current = c
			depth++
			buildCommand2(c, keys[depth:], values)
			return
		}
	}

	if !found {
		for i := depth; i < len(keys); i++ {
			sub := &Command{
				Name:  keys[i],
				Flags: make(map[string]string),
			}
			current.Sub = append(current.Sub, sub)
			current = sub
		}
		for _, value := range values {
			if value[0] == '-' {
				if current.Flags == nil {
					current.Flags = make(map[string]string)
				}
				if _, ok := current.Flags[strings.TrimPrefix(strings.TrimPrefix(value, "-"), "-")]; !ok {
					current.Flags[strings.TrimPrefix(strings.TrimPrefix(value, "-"), "-")] = ""
				}
			} else {
				if current.Args == nil {
					current.Args = []string{}
				}
				if !slices.Contains(current.Args, value[1:]) {
					current.Args = append(current.Args, value[1:])
				}
			}
		}
	}
}

func buildCommand(keys, values []string) Command {
	root := Command{}
	current := &root

	for i := 0; i < len(keys); i++ {
		current.Name = keys[i]
		if i+1 < len(keys) {
			sub := Command{}
			current.Sub = append(current.Sub, &sub)
			current = &sub
		}
	}

	for _, value := range values {
		if value[0] == '-' {
			if current.Flags == nil {
				current.Flags = make(map[string]string)
			}
			if _, ok := current.Flags[strings.TrimPrefix(strings.TrimPrefix(value, "-"), "-")]; !ok {
				current.Flags[strings.TrimPrefix(strings.TrimPrefix(value, "-"), "-")] = ""
			}
		} else {
			if current.Args == nil {
				current.Args = []string{}
			}
			current.Args = append(current.Args, value[1:])
		}
	}

	return root
}

func buildMap(cmd Command) map[string]interface{} {
	yamlConfig := make(map[string]interface{})
	var root interface{}
	buildNode(cmd, &root)
	yamlConfig[cmd.Name] = root
	return yamlConfig
}

func buildNode(cmd Command, yamlNode *interface{}) {
	elements := []interface{}{}

	// Add arguments with '+' prefix
	for _, arg := range cmd.Args {
		elements = append(elements, "+"+arg)
	}

	// Add flags with '--' prefix
	for flag := range cmd.Flags {
		elements = append(elements, "--"+flag)
	}

	// Add subcommands
	for _, sub := range cmd.Sub {
		if len(sub.Sub) == 0 && len(sub.Args) == 0 && len(sub.Flags) == 0 {
			elements = append(elements, sub.Name)
			continue
		}
		subNode := make(map[string]interface{})
		var subYamlNode interface{}
		buildNode(*sub, &subYamlNode)
		subNode[sub.Name] = subYamlNode
		elements = append(elements, subNode)
	}
	if len(elements) > 0 {
		*yamlNode = elements
	} else {
		*yamlNode = nil
	}
}

type Command struct {
	Name     string
	FuncName string
	PkgName  string
	PkgPath  string
	Args     []string
	Flags    map[string]string
	Sub      []*Command
	Parent   *Command
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
					item = strings.TrimPrefix(item, "-")
					item = strings.TrimPrefix(item, "-")
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

func loadConfig() *CLIConfig {
	config := &CLIConfig{Cmd: make(map[string]interface{})}

	var yamlConfig map[string]interface{}
	f := "config.yaml"
	yamlData, err := os.ReadFile(f)
	err = yaml.Unmarshal(yamlData, &yamlConfig)
	if err != nil {
		panic(err)
	}

	config.Cmd = yamlConfig

	return config
}

func saveConfig(config *CLIConfig) {
	data, err := yaml.Marshal(config.Cmd)
	if err != nil {
		fmt.Printf("Error generating YAML: %v\n", err)
		return
	}

	err = ioutil.WriteFile(configFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
	}
}
