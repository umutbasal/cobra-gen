package main

import (
	"os"
	"os/exec"
	"path"
	"strings"

	_ "github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

var modName string

func init() {
	f := "go.mod"
	// read go.mod

	goMod, err := os.ReadFile(f)
	if err != nil {
		panic(err)
	}

	// parse go.mod
	modFile, err := modfile.Parse(f, goMod, nil)
	if err != nil {
		panic(err)
	}

	modName = modFile.Module.Mod.Path
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

func main() {

	if len(os.Args) > 1 {
		generate()
		return
	}
	c := loadConfig()

	commands := parseYaml(c.Cmd)
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
	//exec gofmt -w . in cmd

	cmd := exec.Command("gofmt", "-w", ".")
	cmd.Dir = "cmd"
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	// goimports -w .
	cmd = exec.Command("goimports", "-w", ".")
	cmd.Dir = "cmd"
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	// create examples/main.go
	err = os.MkdirAll("examples", 0755)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("examples/main.go", []byte(exampleTmpl), 0644)
	if err != nil {
		panic(err)
	}
}

type Folder struct {
	Name       string
	SubFolders []*Folder
	Files      []*File
}

type File struct {
	Name          string
	Path          string
	PkgPath       string
	ParentPkg     string
	ParentPkgPath string
	PkgName       string
	RootPkgName   string
	Cmd           *Command
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
		if folder.Name == "cmd" {
			file.Cmd.FuncName = ""
		}
		folder.Name = pkgNaming(folder.Name)
		file.Path = path + "/" + folder.Name + "/" + file.Name
		file.PkgPath = path + "/" + folder.Name
		file.PkgPath = strings.TrimPrefix(file.PkgPath, "./")
		file.ParentPkgPath = strings.TrimPrefix(path, "./")
		if file.Cmd.Parent != nil {
			file.ParentPkg = pkgNaming(file.Cmd.Parent.Name)
		}
		file.RootPkgName = modName
		file.PkgName = folder.Name
		file.Cmd.FuncName = kebabToCamel(file.Cmd.Name)
		file.Cmd.PkgName = pkgNaming(file.Cmd.Name)
		file.Cmd.PkgPath = file.PkgPath
		*files = append(*files, *file)
	}
	for _, sub := range folder.SubFolders {
		printFullPaths(sub, path+"/"+folder.Name, files)
	}
}

func kebabToCamel(s string) string {
	var result string
	for _, part := range strings.Split(s, "-") {
		result += strings.Title(part)
	}
	return result
}

func pkgNaming(s string) string {
	return strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(s, "_", ""), "-", ""))
}
