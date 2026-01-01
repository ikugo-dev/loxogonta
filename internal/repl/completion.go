package repl

import (
	"sort"

	"github.com/c-bata/go-prompt"
)

type CompletionType int

const (
	CompletionType_Keyword CompletionType = iota
	CompletionType_Variable
	CompletionType_Function
)

type completionSymbol struct {
	Name  string
	Kind  CompletionType
	Value string
}

var completions = map[string]completionSymbol{}

func init() {
	keywords := []string{
		"and", "class", "else", "false", "for",
		"fun", "if", "nil", "or", "print",
		"return", "super", "this", "true", "var", "while",
	}

	for _, kw := range keywords {
		completions[kw] = completionSymbol{
			Name: kw,
			Kind: CompletionType_Keyword,
		}
	}
}

func AddCompletion(name string, kind CompletionType, value string) {
	completions[name] = completionSymbol{
		Name:  name,
		Kind:  kind,
		Value: value,
	}
}

func GetCompleter(d prompt.Document) []prompt.Suggest {
	word := d.GetWordBeforeCursor()
	suggestions := []prompt.Suggest{}

	for _, sym := range sortCompletionSymbols() {
		switch sym.Kind {
		case CompletionType_Function:
			suggestions = append(suggestions, prompt.Suggest{
				Text:        sym.Name + "()",
				Description: "function: " + sym.Value,
			})

		case CompletionType_Variable:
			suggestions = append(suggestions, prompt.Suggest{
				Text:        sym.Name + ";",
				Description: "variable: " + sym.Value,
			})

		case CompletionType_Keyword:
			suggestions = append(suggestions, prompt.Suggest{
				Text:        sym.Name + " ",
				Description: "keyword",
			})
		}
	}

	return prompt.FilterHasPrefix(suggestions, word, true)
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
