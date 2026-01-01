package repl

import (
	"sort"
	"strings"
)

type CompletionType int

const (
	CompletionType_Keyword CompletionType = iota
	CompletionType_Variable
	CompletionType_Function
)

type completionSymbol struct {
	Name string
	Kind CompletionType
}

var completions = map[string]completionSymbol{}

func init() {
	keywords := []string{"and", "class", "else", "false", "for", "fun", "if",
		"nil", "or", "print", "return", "super", "this", "true", "var", "while"}

	for _, kw := range keywords {
		completions[kw] = completionSymbol{
			Name: kw,
			Kind: CompletionType_Keyword,
		}
	}
}

func AddCompletion(name string, kind CompletionType, value string) {
	completions[name] = completionSymbol{
		Name: name,
		Kind: kind,
	}
}

func Completer(input string) (c []string) {
	var suggestions []string

	for _, symbol := range sortCompletionSymbols() {
		if !strings.HasPrefix(symbol.Name, input) {
			continue
		}
		switch symbol.Kind {
		case CompletionType_Function:
			suggestions = append(suggestions, symbol.Name+"()")

		case CompletionType_Variable:
			suggestions = append(suggestions, symbol.Name+";")

		case CompletionType_Keyword:
			suggestions = append(suggestions, symbol.Name+" ")
		}
	}
	return suggestions
}

func sortCompletionSymbols() []completionSymbol {
	symbols := make([]completionSymbol, 0, len(completions))
	for _, sym := range completions {
		symbols = append(symbols, sym)
	}

	sort.Slice(symbols, func(i, j int) bool {
		if symbols[i].Kind != symbols[j].Kind {
			return symbols[i].Kind < symbols[j].Kind
		}
		return symbols[i].Name < symbols[j].Name
	})

	return symbols
}
