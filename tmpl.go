package main

import (
	"bytes"
	"text/template"
)

var rootTmpl = `package main

import "{{ .PkgName }}/cmd"
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
