package main

import (
	"fmt"
	"os"
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

var hadError = false

func runFile(filePath string) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Could not read specified file: %s\n%s\n", filePath, err)
		return
	}
	run(string(fileContent))

	if hadError {
		os.Exit(65)
	}
}

func runPrompt() {
	for true {
		fmt.Print("> ")
		var line string
		_, err := fmt.Scanln(&line)
		if err != nil {
			break
		}
		run(line)
		hadError = false
	}
}

func run(source string) {
	fmt.Println(source)

	var scanner Scanner = Scanner{source: source, line: 1}
	var tokens []Token = scanner.scanTokens()
	// For now, just print the tokens.
	for _, token := range tokens {
		fmt.Println("Token:", token)
	}
}
