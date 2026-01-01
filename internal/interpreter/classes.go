package intr

type loxClass struct {
	name     string
	toString func() string
}

func createLoxClass(name string) *loxClass {
	return &loxClass{
		name:     name,
		toString: func() string { return name },
	}
}
