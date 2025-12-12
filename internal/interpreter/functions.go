package intr

import (
	"time"

	"github.com/ikugo-dev/loxogonta/internal/ast"
)

func addNativeFunctions() {
	storage.put("clock", &nativeFunction{
		fnArity: 0,
		fn: func(args []any) any {
			return time.Now()
		},
	})
}

type loxCallable interface {
	call(arguments []any) any
	arity() int
}

type nativeFunction struct {
	fnArity int
	fn      func(args []any) any
}

func (n *nativeFunction) arity() int {
	return n.fnArity
}
func (n *nativeFunction) call(args []any) any {
	return n.fn(args)
}

type loxFunction struct {
	declaration *ast.FunctionStmt
	closure     *environment
}

func (f *loxFunction) arity() int {
	return len(f.declaration.Params)
}
func (f *loxFunction) call(args []any) any {
	// New environment whose parent is the closure
	env := createEnvironmentWithParent(*f.closure)

	// Bind params
	for i, param := range f.declaration.Params {
		env.put(param.Lexeme, args[i])
	}

	// Execute function body
	old := storage
	storage = env

	var result any
	for _, stmt := range f.declaration.Body {
		result = evalStmt(stmt)
		// Handle returns later
	}

	storage = old
	return result
}
