package main

import (
	"flag"
	"log"
)

var outputDir string

func init() {
	flag.StringVar(&outputDir, "o", "", "generate_ast <output directory>")
	flag.Parse()
}

func main() {
	if len(outputDir) == 0 {
		log.Fatalln("Usage: generate_ast <output directory>")
	}

	// 根据文法自动生成表达式的代码
	//err := defineAst(outputDir, "Expr", []string{
	//	"Binary   : left Expr, operator *scanner.Token, right Expr",
	//	"Grouping : expression Expr",
	//	"Literal  : value interface{}",
	//	"Unary    : operator *scanner.Token, right Expr",
	//})
	//
	//if err != nil {
	//	log.Fatalln(err)
	//}

	_ = defineAst(outputDir, "Stmt", []string{
		"ExprStmt : expr Expr",
		"PrintStmt: expr Expr",
	})
}
