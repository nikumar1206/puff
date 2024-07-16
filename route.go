package puff

import (
	"fmt"
	"regexp"
)

type Route struct {
	WebSocket bool
	Protocol  string
	Pattern   string // TODO: remove this. un-necessary
	Path      string
	Handler   func(*Context)
	fullPath  string
	regexp    *regexp.Regexp
	Fields    Field
	// should probably have responses (200 OK followed by json, 400 Invalid Paramaters, etc...)
	// Responses []map[int]Response -> responses likely will look something like this
}

func (r *Route) String() string {
	return fmt.Sprintf("Protocol: %s\nPattern: %s\nPath: %s\n", r.Protocol, r.Pattern, r.Path)
}

// func (r *Route) handleHandlerSchema() {
// 	// Get the type of the handler function
// 	handlerType := reflect.TypeOf(r.Handler)

// 	// Ensure it's a function
// 	if handlerType.Kind() != reflect.Func {
// 		panic("Handler provided MUST be a function.") // likely impossible if case
// 	}

// 	// Extract the input parameter (second parameter)
// 	inputParamType := handlerType.In(1)

// 	// Ensure the input parameter is a pointer to a struct
// 	if inputParamType.Kind() != reflect.Ptr || inputParamType.Elem().Kind() != reflect.Struct {
// 		panic("Input parameter MUST be a pointer to a struct.")
// 	}

// 	// Get the struct type
// 	structType := inputParamType.Elem()

// 	r.Schema = structType

// 	fields := make(map[string]Param)
// 	for i := 0; i < structType.NumField(); i++ {
// 		field := structType.Field(i)

// 		field_type_name := field.Type.Name()
// 		if field_type_name != "string" && field_type_name != "int" && field_type_name != "bool" {
// 			panic("Field type of input struct must be string, int, or bool.")
// 		}

// 		description := field.Tag.Get("description")
// 		kind := field.Tag.Get("kind")

// 		if kind != "HeaderParam" && kind != "PathParam" && kind != "QueryParam" && kind != "BodyParam" {
// 			panic("Kind must be specified and must be HeaderParam, PathParam, QueryParam, or BodyParam.")
// 		}

// 		fields[field.Name] = Param{
// 			Description: description,
// 			Kind:        kind,
// 			Type:        field_type_name,
// 		}
// 	}
// 	r.Parameters = fields
// }
