package main

import (
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/ikugo-dev/loxogonta/internal/errors"
	"github.com/ikugo-dev/loxogonta/internal/interpreter"
	"github.com/ikugo-dev/loxogonta/internal/parser"
	"github.com/ikugo-dev/loxogonta/internal/repl"
	"github.com/ikugo-dev/loxogonta/internal/scanner"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: ./main [script]")
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

func run(source string) any {
	tokens := scn.ScanSource(source)
	statements := prs.ParseTokens(tokens)
	if errors.HadError {
		return nil
	}
	return intr.Interpret(statements)
}

func runPrompt() {
	p := prompt.New(
		func(line string) {
			if lastExpr := run(line); lastExpr != nil {
				fmt.Println("--> ", lastExpr)
			}
			errors.HadError = false
		},
		repl.GetCompleter,
		prompt.OptionPrefix("> "),
		prompt.OptionHistory([]string{}),
	)

	p.Run()
}
