package puff

type Field struct {
	Name        string
	Description string
	Body        map[string]any
	// by default not required. unless specified
	QueryParams map[string]any
	// by default required. unless specified
	PathParams string
}

type Query struct {
	foo string
	bar string
}
