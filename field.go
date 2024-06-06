package puff

type Param interface {
	GetKind() string        //Query, Path, Header
	GetDescription() string //for the OpenAPI
	GetValue() interface{}
	SetValue()
}

type PathParam[T any] struct {
	Description string
	Value       T
}

func (p PathParam[T]) GetKind() string {
	return "Path"
}
func (p PathParam[T]) GetDescription() string {
	return p.Description
}
func (p PathParam[T]) GetValue() interface{} {
	return p.Value
}
func (p *PathParam[T]) SetValue(v T) {
	p.Value = v
}

type QueryParam[T any] struct {
	Description string
	Value       T
}

func (p QueryParam[T]) GetKind() string {
	return "Query"
}
func (p QueryParam[T]) GetDescription() string {
	return p.Description
}
func (p QueryParam[T]) GetValue() interface{} {
	return p.Value
}
func (p *QueryParam[T]) SetValue(v T) {
	p.Value = v
}

type HeaderParam[T any] struct {
	Description string
	Value       T
}

func (p HeaderParam[T]) GetKind() string {
	return "Header"
}
func (p HeaderParam[T]) GetDescription() string {
	return p.Description
}
func (p HeaderParam[T]) GetValue() interface{} {
	return p.Value
}
func (p *HeaderParam[T]) SetValue(v T) {
	p.Value = v
}

//EXAMPLE FIELDS SCHEMA PASSED IN:

// type HelloWorld struct {
// 	Globe HeaderParam[string]
// }

//globe will be the name, header param will be the kind, type of the value it should get
//is string, however Description remains empty for now
