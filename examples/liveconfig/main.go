package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"slices"

	"gopkg.in/yaml.v3"
)

type Command map[string]interface{}

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

	current := generateNestedMap(keys, values)

	// shitcode to reset the refs
	var source map[string]interface{}
	sourcedata, err := yaml.Marshal(current)
	if err != nil {
		fmt.Printf("Error marshalling current: %v\n", err)
		return
	}
	err = yaml.Unmarshal(sourcedata, &source)

	err = mergeMaps(config.Cmd, source)
	if err != nil {
		fmt.Printf("Error merging maps: %v\n", err)
		return
	}

	c, _ := yaml.Marshal(config)
	fmt.Println(string(c))

	saveConfig(config)
}

func mergeMaps(dst, src map[string]interface{}) error {
	for srcKey, srcValue := range src {
		if srcValueAsMap, ok := srcValue.(map[string]interface{}); ok { // handle maps
			if dstValue, ok := dst[srcKey]; ok {
				if dstValueAsMap, ok := dstValue.(map[string]interface{}); ok {
					err := mergeMaps(dstValueAsMap, srcValueAsMap)
					if err != nil {
						return err
					}
					continue
				}
				// If dstValue is not a map, create a new map and continue merging
				dst[srcKey] = make(map[string]interface{})
				err := mergeMaps(dst[srcKey].(map[string]interface{}), srcValueAsMap)
				if err != nil {
					return err
				}
			} else {
				dst[srcKey] = make(map[string]interface{})
				err := mergeMaps(dst[srcKey].(map[string]interface{}), srcValueAsMap)
				if err != nil {
					return err
				}
			}
		} else if srcValueAsSlice, ok := srcValue.([]interface{}); ok { // handle slices
			if dstValue, ok := dst[srcKey]; ok {
				if dstValueAsSlice, ok := dstValue.([]interface{}); ok {
					for _, srcValueAsSliceElement := range srcValueAsSlice {
						if !slices.Contains(dstValueAsSlice, srcValueAsSliceElement) {
							dstValueAsSlice = append(dstValueAsSlice, srcValueAsSliceElement)
						}
					}
					dst[srcKey] = dstValueAsSlice
					continue
				}
				// If dstValue is not a slice, replace with srcValueAsSlice
				dst[srcKey] = srcValueAsSlice
			} else {
				dst[srcKey] = srcValueAsSlice
			}
		} else { // handle other types
			dst[srcKey] = srcValue
		}
	}
	return nil
}

func generateNestedMap(keys []string, values []string) map[string]interface{} {
	output := make(map[string]interface{})

	if len(keys) == 0 {
		return output
	}

	current := output
	for i, key := range keys {
		if i == len(keys)-1 {
			current[key] = values
		} else {
			next := make(map[string]interface{})
			current[key] = next
			current = next
		}
	}

	return output
}
func loadConfig() *CLIConfig {
	config := &CLIConfig{Cmd: make(map[string]interface{})}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return config
		}
		fmt.Printf("Error reading config file: %v\n", err)
		return config
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
	}

	return config
}

func saveConfig(config *CLIConfig) {
	data, err := yaml.Marshal(config)
	if err != nil {
		fmt.Printf("Error generating YAML: %v\n", err)
		return
	}

	err = ioutil.WriteFile(configFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
	}
}
