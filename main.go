package main

import (
	"bufio"
	"fmt"
	"os"

	"ksm/interpreter"
	"ksm/lexer"
	"ksm/parser"
)

func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	i := interpreter.NewInterpreter()

	for {
		fmt.Print(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}

		l := lexer.NewLexer(line)
		p := parser.NewParser(l)

		ast, err := p.Parse()
		if err != nil {
			fmt.Printf("Parsing error: %v\n", err)
			continue
		}

		err = i.Interpret(ast)
		if err != nil {
			fmt.Printf("Interpretation error: %v\n", err)
		}
	}
}

func main() {
	fmt.Println("Welcome to the KSM REPL!")
	fmt.Println("Type 'exit' to quit.")
	repl()
}
