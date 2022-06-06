package main

import (
	"flag"
)

var source string

func init() {
	flag.StringVar(&source, "s", "", "Lox source code file path")
	flag.Parse()
}

func main() {
	runApp(source)
}
