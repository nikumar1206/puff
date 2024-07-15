package puff

import (
	"fmt"
	"reflect"
)

type Field struct {
	// Description is the Route's description. To be reflected in OpenAPI spec
	Description string
	Body        map[string]any
	// by default not required. unless specified
	QueryParams map[string]any
	// by default required. unless specified
	PathParams map[string]reflect.Kind
	Responses  map[int]Response
	Validators []func() bool
}

type Query struct {
	foo string
	bar string
}

func (f *Field) ValidateIncomingAttribute(foo reflect.StructField, v any) error {
	fmt.Println("found the following provided", foo.Type, v)

	var err error
	if err != nil {
		err = fmt.Errorf("So wrong")
	}
	return err
}
