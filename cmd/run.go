package main

import (
	"GLox/interpreter"
	le "GLox/loxerror"
	"GLox/parser"
	"GLox/resolver"
	"GLox/scanner"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func runApp(source string) {
	if source != "" {
		runFile(source)
	} else {
		fmt.Println("please specific the source code file.")
	}
}

func runFile(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	sc := string(bytes)

	run(sc)
}

func run(sc string) {
	s := scanner.NewScanner(sc)
	tokens := s.ScanTokens()

	p := parser.NewParser(tokens)
	stmts := p.Parse()
	if le.HadError {
		os.Exit(-1)
	}

	i := interpreter.NewInterpreter()
	r := resolver.NewResolver(i)
	r.ResolveStmt(stmts...)
	if le.HadResolveError {
		os.Exit(-2)
	}

	err := i.Interpret(stmts)
	if err != nil {
		fatal(err.Error(), 0)
	}
}

func fatal(msg string, signal int) {
	fmt.Println(msg)
	os.Exit(signal)
}
