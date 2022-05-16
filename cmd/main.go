package main

import (
	"flag"
)

var source string

func init() {
	flag.StringVar(&source, "s", "", "Lox source code file")
	flag.Parse()
}

func main() {
	if source != "" {
		runFile(source)
	} else {
		runPrompt()
	}
}
