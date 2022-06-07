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

	if le.HadRuntimeError {
		os.Exit(-2)
	}
}

func run(sc string) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err.(error).Error())
		}
	}()

	s := scanner.NewScanner(sc)
	tokens := s.ScanTokens()

	p := parser.NewParser(tokens)
	stmts, err := p.Parse()
	if err != nil {
		fatal(err.Error(), -1)
	}

	i := interpreter.NewInterpreter()
	r := resolver.NewResolver(i)
	r.ResolveStmt(stmts...)

	if le.HadError {
		return
	}

	i.Interpret(stmts)
}

func fatal(msg string, n int) {
	fmt.Println(msg)
	os.Exit(n)
}
