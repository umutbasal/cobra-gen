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

func generate() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("No commands provided.")
		return
	}

	config := loadConfig()

	commands, params := partitionArgs(args)
	cmd := parseYaml(config.Cmd)

	updateCommands(&cmd, commands, params)
	config.Cmd = buildMap(cmd)

	saveConfig(config)
}

// Helper function to partition arguments into commands and params
func partitionArgs(args []string) (commands []string, params []string) {
	for _, arg := range args {
		if arg[0] == '-' || arg[0] == '+' {
			params = append(params, arg)
		} else {
			commands = append(commands, arg)
		}
	}
	return commands, params
}

func updateCommands(root *Command, commands, params []string) {
	current := root
	depth := 0

	for {
		// If we have processed all command, add the params to the current command
		if depth >= len(commands) {
			addParams(current, params)
			return
		}

		// Search for an last subcommand matching the current command
		found := false
		for _, sub := range current.Sub {
			if sub.Name == commands[depth] {
				current = sub
				found = true
				depth++
				break
			}
		}

		// If no matching subcommand was found, create a new one
		if !found {
			for i := depth; i < len(commands); i++ {
				newSub := &Command{
					Name:  commands[i],
					Flags: make(map[string]string),
				}
				current.Sub = append(current.Sub, newSub)
				current = newSub
			}
			addParams(current, params)
			return
		}
	}
}

// Helper function to add params to the command
func addParams(cmd *Command, params []string) {
	for _, param := range params {
		if param[0] == '-' {
			flagName := strings.TrimPrefix(strings.TrimPrefix(param, "-"), "-")
			if cmd.Flags == nil {
				cmd.Flags = make(map[string]string)
			}
			if _, exists := cmd.Flags[flagName]; !exists {
				cmd.Flags[flagName] = ""
			}
		} else {
			arg := param[1:]
			if cmd.Args == nil {
				cmd.Args = []string{}
			}
			if !slices.Contains(cmd.Args, arg) {
				cmd.Args = append(cmd.Args, arg)
			}
		}
	}
}

func buildMap(cmd Command) (m map[string]interface{}) {
	m = make(map[string]interface{})
	m[cmd.Name] = buildNode(cmd)
	return m
}

func buildNode(cmd Command) interface{} {
	var elements []interface{}

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
		} else {
			subNode := make(map[string]interface{})
			subNode[sub.Name] = buildNode(*sub)
			elements = append(elements, subNode)
		}
	}

	if len(elements) > 0 {
		return elements
	}
	return nil
}

func loadConfig() *CLIConfig {
	config := &CLIConfig{Cmd: make(map[string]interface{})}

	var yamlConfig map[string]interface{}
	yamlData, err := os.ReadFile(configFile)
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

	fmt.Println("Config file updated.")
	fmt.Println(string(data))
}
