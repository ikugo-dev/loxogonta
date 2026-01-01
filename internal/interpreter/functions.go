package intr

import (
	"time"

	"github.com/ikugo-dev/loxogonta/internal/ast"
)

func addNativeFunctions(storage *environment) {
	storage.put("clock", &nativeFunction{
		fnArity: 0,
		fn: func(args []any) any {
			return time.Now()
		},
	})
}

type loxCallable interface {
	call(storage *environment, arguments []any) any
	arity() int
}

type nativeFunction struct {
	fnArity int
	fn      func(args []any) any
}

func (n *nativeFunction) arity() int {
	return n.fnArity
}
func (n *nativeFunction) call(_ *environment, args []any) any {
	return n.fn(args)
}

type loxFunction struct {
	declaration *ast.FunctionStmt
	closure     *environment
}

func (f *loxFunction) arity() int {
	return len(f.declaration.Params)
}
func (f *loxFunction) call(storage *environment, args []any) any {
	oldStorage := storage
	storage = createEnvironmentWithParent(f.closure)
	for i, param := range f.declaration.Params {
		storage.put(param.Lexeme, args[i])
	}

	var result any
	for _, stmt := range f.declaration.Body {
		result = evalStmt(storage, stmt)
	}

	storage = oldStorage
	return result
}
