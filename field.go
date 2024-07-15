package puff

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// Param interface
type Param struct {
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
	foo  string
	bar  string
	Kind string
	Type string
}
type NoParams struct{}

func NewParametersInterface(request *http.Request, schema reflect.Type, params map[string]Param) (*interface{}, error) {
	if len(params) == 0 {
		return nil, nil
	}

	newStruct := reflect.New(schema).Elem()

	for field_name, param := range params {
		schema_field := newStruct.FieldByName(field_name)

		if !schema_field.IsValid() { //should never hit this condition, but specified
			return nil, fmt.Errorf("Unexpected parameter: %s", field_name)
		}

		reqBody, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, fmt.Errorf("Unexpected error: %s", err.Error())
		}
		var body map[string]any
		if string(reqBody) != "" {
			err := json.Unmarshal(reqBody, &body)
			if err != nil {
				return nil, fmt.Errorf("Body formatted incorrectly. Error: %s.", err.Error())
			}
		}

		var value any
		switch param.Type {
		case "QueryParam":
			value = request.URL.Query().Get(field_name)
		case "HeaderParam":
			value = request.Header.Get(field_name)
		case "BodyParam":
			if len(body) == 0 {
				value = ""
				break
			}
			val, ok := body[field_name]
			if !ok {
				value = ""
				break
			}
			value = val
		case "PathParam":
			value = request.PathValue(field_name)
		}
		switch typedValue := value.(type) {
		case string:
			if param.Type != "string" {
				return nil, fmt.Errorf("Parameter %s must be a %s.", field_name, param.Type)
			}
			schema_field.SetString(typedValue)
			if value == "" {
				return nil, fmt.Errorf("Expected parameter but not provided: %s", field_name)
			}
		case int:
			if param.Type != "int" {
				return nil, fmt.Errorf("Parameter %s must be a %s.", field_name, param.Type)
			}
			schema_field.SetInt(int64(typedValue))
		case bool:
			if param.Type != "bool" {
				return nil, fmt.Errorf("Parameter %s must be a %s.", field_name, param.Type)
			}
			schema_field.SetBool(typedValue)
		default:
			return nil, fmt.Errorf("Parameter %s must be a %s.", field_name, param.Type)
		}
	}
	ns, ok := newStruct.Interface().(interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected error occurred.")
	}
	return &ns, nil
}
