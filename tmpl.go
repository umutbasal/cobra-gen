package main

import (
	"bytes"
	"text/template"
)

var rootTmpl = `package {{ .PkgName }}

{{ $r := .RootPkgName }}
import (
	"fmt"
	"github.com/spf13/cobra"
	{{ if ne .PkgName .ParentPkg}}
	{{ if .Cmd.Sub}}
	{{ range .Cmd.Sub }}
	{{ if .Sub}}
	"{{ $r }}/{{ .PkgPath }}"
	{{ end }}
	{{ end }}
	{{ end }}
	{{ end }}
)
{{ $base := .Cmd }}
var {{ if ne .Cmd.FuncName "Cmd"}}{{ .Cmd.FuncName }}{{ end }}Cmd = &cobra.Command{
	Use:   "{{ .Cmd.Name }}",
	Short: "{{ .Cmd.Name }} ...",
	{{ if .Cmd.Args }}
	Args: cobra.MatchAll(cobra.ExactArgs({{ len .Cmd.Args }}), cobra.OnlyValidArgs),
	{{ end }}
	Run: func(cmd *cobra.Command, args []string) {
		{{if .Cmd.Args }}
		{{ range $k, $arg := .Cmd.Args }}
		{{ $arg }} := args[{{ $k }}]
		fmt.Printf("{{ $arg }}: %v\n", {{ $arg }})
		{{ end }}
		{{ end }}
		// TODO: Implement command
	},
}

{{if .Cmd.Flags}}
var test{{ .Cmd.FuncName }} string
{{ end }}

func init() {
	{{$pkg := .PkgName}}
	{{$parent := .ParentPkg}}
	{{$cmd := .Cmd}}
	// Flags
	{{ if .Cmd.Flags}}
	{{ range $key, $flag := .Cmd.Flags }}
	// {{ $key }} : {{ $flag }}
	{{ if eq $pkg $parent }}{{ if $cmd }}{{ $cmd.Parent.FuncName }}Cmd.Flags().StringVarP(&test{{ $cmd.FuncName }}, "{{ $key }}", "{{$key}}", "", "desc"){{end}}
	{{ else }}
	{{ if $cmd }}{{ if ne $cmd.FuncName "Cmd"}}{{ $cmd.FuncName }}{{ end }}Cmd.PersistentFlags().StringVarP(&test{{ $cmd.FuncName }}, "{{ $key }}", "{{$key}}", "", "desc") {{end}}
	{{end}}
	{{ end }}
	{{ end }}
	
	{{ if eq .PkgName .ParentPkg}}{{ if .Cmd.Parent }}{{ .Cmd.Parent.FuncName }}Cmd.AddCommand({{ .Cmd.FuncName }}Cmd)
	{{ end }}{{ end }}
	{{ if ne .PkgName .ParentPkg}}
	{{ if .Cmd.Sub}}
	{{ range .Cmd.Sub }}
	{{ if .Sub}}
	{{ if ne $base.FuncName "Cmd"}}{{ $base.FuncName }}{{ end }}Cmd.AddCommand({{.PkgName}}.{{ .FuncName }}Cmd)
	{{ end }}
	{{ end }}
	{{ end }}
	{{ end }}
}
`

var exampleTmpl = `package main

import (
	"github.com/umutbasal/cobra-gen/cmd"
)


func main(){
	cmd.Cmd.Execute()
}
`

func execTmpl(f *File) []byte {
	tmpl, err := template.New("main").Parse(rootTmpl)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, f)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}
