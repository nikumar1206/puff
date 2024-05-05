package puff

type Field struct {
	Name        string
	Schema      interface{}
	Description string
}

type QueryParam struct {
	*Field
}

type PathParam struct {
	*Field
}

type Header struct {
	*Field
}
