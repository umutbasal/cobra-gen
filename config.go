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

	keys := []string{}
	values := []string{}
	for _, arg := range args {
		if arg[0] == '-' || arg[0] == '+' {
			values = append(values, arg)
		} else {
			keys = append(keys, arg)
		}
	}

	cmd := parseYaml(config.Cmd)

	updateCommands(&cmd, keys, values)
	resm := buildMap(cmd)

	config.Cmd = resm

	saveConfig(config)
}

func updateCommands(root *Command, keys, values []string) {
	current := root
	depth := 0

	for {
		// If we have processed all keys, add the values
		if depth >= len(keys) {
			addValues(current, values)
			return
		}

		// Search for an existing subcommand matching the current key
		found := false
		for _, sub := range current.Sub {
			if sub.Name == keys[depth] {
				current = sub
				found = true
				depth++
				break
			}
		}

		// If no matching subcommand was found, create a new one
		if !found {
			for i := depth; i < len(keys); i++ {
				newSub := &Command{
					Name:  keys[i],
					Flags: make(map[string]string),
				}
				current.Sub = append(current.Sub, newSub)
				current = newSub
			}
			addValues(current, values)
			return
		}
	}
}

// Helper function to add values to the command
func addValues(cmd *Command, values []string) {
	for _, value := range values {
		if value[0] == '-' {
			flagName := strings.TrimPrefix(strings.TrimPrefix(value, "-"), "-")
			if cmd.Flags == nil {
				cmd.Flags = make(map[string]string)
			}
			if _, exists := cmd.Flags[flagName]; !exists {
				cmd.Flags[flagName] = ""
			}
		} else {
			arg := value[1:]
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
