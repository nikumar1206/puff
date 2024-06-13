package puff

type Field struct {
	Description string
	Body        map[string]any
	// by default not required. unless specified
	QueryParams map[string]any
	// by default required. unless specified
	PathParams string
	Responses  map[int]Response
}

type Query struct {
	foo string
	bar string
}
