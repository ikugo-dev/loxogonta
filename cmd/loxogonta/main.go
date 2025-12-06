package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ikugo-dev/loxogonta/internal/errors"
	intr "github.com/ikugo-dev/loxogonta/internal/interpreter"
	"github.com/ikugo-dev/loxogonta/internal/parser"
	"github.com/ikugo-dev/loxogonta/internal/scanner"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: jlox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(filePath string) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Could not read specified file: %s\n%s\n", filePath, err)
		return
	}
	run(string(fileContent))

	if errors.HadError || errors.HadParseError {
		os.Exit(65)
	}
	if errors.HadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt() {
	in := bufio.NewReader(os.Stdin)
	for true {
		fmt.Print("> ")
		var line string
		line, err := in.ReadString('\n')
		if err != nil {
			break
		}
		run(line)
		errors.HadError = false
	}
}

func run(source string) {
	scanner := scn.NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := prs.NewParser(tokens)
	statements := parser.Parse()
	if errors.HadError {
		return
	}
	for _, statement := range statements {
		fmt.Printf("statement: %v\n", statement)
	}
	intr.Interpret(statements)
}
