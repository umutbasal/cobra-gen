package main

import (
	"fmt"
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
	Args     []string
	Flags    map[string]string
	File     *File
	Sub      []*Command
	Parent   *Command
}

func parseYaml(m map[string]interface{}) Command {
	if len(m) != 1 {
		panic("Root command must be only one")
	}
	rootStr := ""
	for key := range m {
		rootStr = key
	}
	root := Command{Name: rootStr, Flags: make(map[string]string)}
	parseNode(rootStr, m[rootStr], &root)
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
					flag := strings.TrimPrefix(strings.TrimPrefix(item, "-"), "-")
					parent.Flags[flag] = ""
				} else {
					sub := &Command{Name: item, Flags: make(map[string]string), Parent: parent}
					parent.Sub = append(parent.Sub, sub)
				}
			case map[interface{}]interface{}:
				parseMap(item, parent)
			}
		}
	case map[interface{}]interface{}:
		parseMap(v, parent)
	}
}

func parseMap(m map[interface{}]interface{}, parent *Command) {
	for key, value := range m {
		sub := &Command{Name: key.(string), Flags: make(map[string]string), Parent: parent}
		parseNode(key.(string), value, sub)
		parent.Sub = append(parent.Sub, sub)
	}
}

func main() {
	if len(os.Args) > 1 {
		updateConfig()
		return
	}

	formatCode()

	c := loadConfig()
	commands := parseYaml(c.Cmd)
	// if cmd exists throw error
	if _, err := os.Stat("cmd"); err == nil {
		panic("cmd directory already exists")
	}
	folders := &Folder{Name: "cmd"}

	structureFolders(commands, 0, folders)
	var files []File
	fillForTemplate(folders, ".", &files)

	createFilesAndDirectories(files)
	formatCode()
	createExamplesMain()
}

func createFilesAndDirectories(files []File) {
	for _, f := range files {
		fmt.Println(f.Path)

		if err := os.MkdirAll(path.Dir(f.Path), 0755); err != nil {
			panic(err)
		}

		if err := os.WriteFile(f.Path, execTmpl(&f), 0644); err != nil {
			panic(err)
		}
	}
}

func formatCode() {

	var r bool
	// check if gofmt exists
	if !commandExists("gofmt") {
		fmt.Println("gofmt not found, install it")
		r = true
	}

	// check if goimports exists
	if !commandExists("goimports") {
		fmt.Println("goimports not found, install it")
		r = true
	}
	if r {
		panic("")
	}

	runCommand("gofmt", "cmd")
	//go install golang.org/x/tools/cmd/goimports@latest
	runCommand("goimports", "cmd")
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func runCommand(name, dir string) {
	cmd := exec.Command(name, "-w", ".")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running %s please install it. (err: %v)\n", name, err)
		panic(err)
	}
}

func createExamplesMain() {
	if err := os.MkdirAll("examples/cobra-gen", 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile("examples/cobra-gen/main.go", []byte(exampleTmpl(modName)), 0644); err != nil {
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
		result.Files = append(result.Files, &File{Name: "cmd.go", Cmd: &cmd})
	}

	for _, sub := range cmd.Sub {
		file := &File{Name: sub.Name + ".go", Cmd: sub}
		if len(sub.Sub) > 0 {
			folder := &Folder{Name: sub.Name}
			result.SubFolders = append(result.SubFolders, folder)
			structureFolders(*sub, level+1, folder)
			folder.Files = append(folder.Files, file)
		} else {
			result.Files = append(result.Files, file)
		}
	}
}

func fillForTemplate(folder *Folder, path string, files *[]File) {
	folder.Name = pkgNaming(folder.Name)
	for _, file := range folder.Files {
		modifyFile(file, folder, path)
		*files = append(*files, *file)
	}
	for _, sub := range folder.SubFolders {
		fillForTemplate(sub, path+"/"+folder.Name, files)
	}
}

func modifyFile(file *File, folder *Folder, path string) {
	if folder.Name == "cmd" {
		file.Cmd.FuncName = ""
	}
	filePath := path + "/" + folder.Name
	file.Path = filePath + "/" + file.Name
	file.PkgPath = strings.TrimPrefix(filePath, "./")
	file.ParentPkgPath = strings.TrimPrefix(path, "./")
	if file.Cmd.Parent != nil {
		file.ParentPkg = pkgNaming(file.Cmd.Parent.Name)
	}
	file.RootPkgName = modName
	file.PkgName = folder.Name
	file.Cmd.FuncName = kebabToCamel(file.Cmd.Name)
	file.Cmd.File = file
}

func kebabToCamel(s string) string {
	parts := strings.Split(s, "-")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}

func pkgNaming(s string) string {
	return strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(s, "_", ""), "-", ""))
}
