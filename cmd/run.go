package main

import (
	"GLox/interpreter"
	le "GLox/lerror"
	"GLox/parser"
	"GLox/resolver"
	"GLox/scanner"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func runFile(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	sc := string(bytes)

	run(sc)
	if le.HadError {
		os.Exit(-1)
	}

	if le.HadRuntimeError {
		os.Exit(-2)
	}
}

func runPrompt() {
	for {
		fmt.Print("> ")
		line, _, err := bufio.NewReader(os.Stdin).ReadLine()
		if err != nil {
			log.Println(err)
		}

		if len(line) == 0 {
			break
		}

		run(string(line))
		le.HadError = false
	}
}

func run(sc string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err.(error).Error())
		}
	}()

	s := scanner.NewScanner(sc)
	tokens := s.ScanTokens()

	p := parser.NewParser(tokens)
	stmts := p.Parse()

	// 如果有语法错误直接返回，不把错误的信息带给解释器
	if le.HadError {
		return
	}

	i := interpreter.NewInterpreter()
	r := resolver.NewResolver(i)

	r.ResolveStmt(stmts...)
	i.Interpret(stmts)

	//fmt.Println(new(parser.Printer).Print(expr))

	//for _, token := range tokens {
	//	fmt.Println(token)
	//}
}
