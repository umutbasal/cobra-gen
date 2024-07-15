package main

import (
	"os"
	"path"
	"strings"

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
	folders := &Folder{
		Name: "cmd",
	}
	structureFolders(commands, 0, folders)
	var files []File
	printFullPaths(folders, ".", &files)
	for _, f := range files {
		println(f.Path)

		err := os.MkdirAll(path.Dir(f.Path), 0755)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(f.Path, execTmpl(&f), 0644)
		if err != nil {
			panic(err)
		}
	}
}

type Folder struct {
	Name       string
	SubFolders []*Folder
	Files      []*File
}

type File struct {
	Name    string
	Path    string
	PkgPath string
	PkgName string
	Cmd     *Command
}

func structureFolders(cmd Command, level int, result *Folder) {
	if len(cmd.Sub) == 0 {
		return
	}
	if level == 0 {
		file := &File{
			Name: "cmd.go",
			Cmd:  &cmd,
		}
		result.Files = append(result.Files, file)
	}
	for _, sub := range cmd.Sub {
		if len(sub.Sub) > 0 {
			folder := &Folder{
				Name: sub.Name,
			}
			result.SubFolders = append(result.SubFolders, folder)
			structureFolders(*sub, level+1, folder)
			file := &File{
				Name: sub.Name + ".go",
				Cmd:  sub,
			}
			folder.Files = append(folder.Files, file)
		} else {
			file := &File{
				Name: sub.Name + ".go",
				Cmd:  sub,
			}
			result.Files = append(result.Files, file)
		}
	}
}

func printFullPaths(folder *Folder, path string, files *[]File) {
	for _, file := range folder.Files {
		file.Path = path + "/" + folder.Name + "/" + file.Name
		file.PkgPath = path + "/" + folder.Name
		file.PkgPath = strings.TrimPrefix(file.PkgPath, "./")
		file.PkgName = "github.com/umutbasal/cobra-gen"
		*files = append(*files, *file)
	}
	for _, sub := range folder.SubFolders {
		printFullPaths(sub, path+"/"+folder.Name, files)
	}
}
