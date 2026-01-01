package main

import (
	"fmt"
	"os"

	"github.com/ikugo-dev/loxogonta/internal/errors"
	"github.com/ikugo-dev/loxogonta/internal/interpreter"
	"github.com/ikugo-dev/loxogonta/internal/parser"
	"github.com/ikugo-dev/loxogonta/internal/repl"
	"github.com/ikugo-dev/loxogonta/internal/scanner"
	"github.com/peterh/liner"
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
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetCompleter(repl.Completer)

	for true {
		if input, err := line.Prompt("> "); err == nil {
			line.AppendHistory(input)
			if lastExpr := run(input); lastExpr != nil {
				fmt.Println("--> ", lastExpr)
			}
			errors.HadError = false
		} else if err == liner.ErrPromptAborted {
			fmt.Print("Aborted")
			break
		} else {
			fmt.Print("Error reading line: ", err)
		}
	}
}
