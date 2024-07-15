package main

func MainTemplate() []byte {
	return []byte(`
package main

import "{{ .PkgName }}/cmd"

func main() {
	cmd.Execute()
}
`)
}
