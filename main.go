package main

import (
	"flag"
	"os"
	"text/template"
)

//go:generate go run main.go -name=A
//go:generate go run main.go -name=B
//go:generate go run main.go -name=C

var name = flag.String("name", "Testa", "name of struct")

var code =
`package main

type Struct{{.}} struct {}

func (s *Struct{{.}} ) Validate() bool {
	return true
}
`

func main() {
	flag.Parse()
	file, _ := os.Create(*name + ".go")
	defer file.Close()

	tmpl, _ := template.New("test").Parse(code)
	tmpl.Execute(file, *name)
}