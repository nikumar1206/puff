package puff

import (
	"reflect"
)

// Field defines required fields to be inputted.
// FIXME: This is not how we are going to do this.
type Field struct {
	Description string
	Body        map[string]any
	// by default not required. unless specified
	QueryParams map[string]any
	// by default required. unless specified
	PathParams map[string]reflect.Kind
	Responses  map[int]Response
	// Validators []func() bool
}
