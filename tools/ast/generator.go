package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

func defineAst(outputDir string, baseName string, types []string) error {
	f := outputDir + "/" + strings.ToLower(baseName) + ".go"
	var buffer bytes.Buffer
	buffer.WriteString("package parser\n\n")
	buffer.WriteString("import \"GLox/scanner\"\n\n")
	buffer.WriteString("type " + baseName + " interface" + " {\n\tAccept(visitor Visitor) interface{}\n}\n\n")

	for _, t := range types {
		ts := strings.Split(t, ":")
		name := strings.TrimSpace(ts[0])
		fieldList := strings.TrimSpace(ts[1])
		defineType(&buffer, name, fieldList)
		defineVisitor(&buffer, name)
		buffer.WriteString("\n")
	}

	return ioutil.WriteFile(f, buffer.Bytes(), 0644)
}

func defineType(buffer *bytes.Buffer, typeName string, fieldList string) {
	buffer.WriteString("type " + typeName + " struct" + " { \n")

	var filedName []string
	fields := strings.Split(fieldList, ", ")
	for _, filed := range fields {
		filedName = append(filedName, strings.Split(filed, " ")[0])
		buffer.WriteString("\t" + filed + "\n")
	}
	buffer.WriteString("}\n\n")
	// 构造方法
	buffer.WriteString("func New" + typeName + "(" + fieldList + ")" + " *" + typeName + " {\n")
	buffer.WriteString(fmt.Sprintf("\treturn &%s{", typeName))
	buffer.WriteString(strings.Join(filedName, ", "))
	buffer.WriteString("}\n}\n\n")
}

func defineVisitor(buffer *bytes.Buffer, typeName string) {
	receiver := string(typeName[0] + 32)
	buffer.WriteString(fmt.Sprintf("func (%s *%s) accept(visitor Visitor) interface{} {\n", receiver, typeName))
	buffer.WriteString(fmt.Sprintf("\treturn visitor.Visit%s(%s)\n}\n\n", typeName, receiver))
}
