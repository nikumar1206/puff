package puff

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// FIXME: allow for example values

// isValidKind takes in specified_kind and returns
// if it is a supported and valid kind
func isValidKind(specified_kind string) bool {
	// https://swagger.io/specification/#:~:text=the%20in%20property.-,in,query%22%2C%20%22header%22%2C%20%22path%22%20or%20%22cookie%22.,-description
	return specified_kind == "header" ||
		specified_kind == "path" ||
		specified_kind == "query" ||
		specified_kind == "cookie" ||
		specified_kind == "body" ||
		specified_kind == "formdata"
}

// handleParam takes the value as recieved, returns an error if the value
// is empty AND required.
func handleParam(value string, param Parameter) (string, error) {
	ok := !(value == "")
	if !ok && param.Required {
		return "", fmt.Errorf("required %s param %s not provided.", param.In, param.Name)
	}
	return value, nil
}

// validate validates the input string against the type to ensure with options
// from Parameter.
func validate(input map[string]any, schemaType reflect.Type) (bool, error) {
	expectedNotFoundKeys := map[string]bool{}
	fields := map[string]reflect.StructField{}
	for i := range schemaType.NumField() {
		field := schemaType.Field(i)
		name := field.Name
		nameTag := field.Tag.Get("name")
		jsonTag := field.Tag.Get("json")
		s := strings.Split(jsonTag, ",")
		jsonTagName := s[0]
		if jsonTagName != "" {
			name = jsonTag
		}
		if nameTag != "" { // name takes priority over json
			name = nameTag
		}
		fields[name] = field
		b, _ := resolveBool(field.Tag.Get("required"), true)
		expectedNotFoundKeys[name] = b
	}
	for k, v := range input {
		required, ok := expectedNotFoundKeys[k]
		if !ok {
			return false, UnexpectedJSONKey(k)
		} else {
			delete(expectedNotFoundKeys, k)
		}
		f, _ := fields[k] //cannot error
		ft := f.Type
		p := ft.Kind() == reflect.Pointer
		tr := reflect.TypeOf(v)
		if tr == nil {
			if required && !p {
				return false, BadFieldType(k, "nil", ft.Kind().String())
			} else {
				continue
			}
		}
		t := tr.Kind()
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
			// p = true
		}
		switch t {
		case reflect.String, reflect.Bool:
			if ft.Kind() != t {
				return false, BadFieldType(k, t.String(), ft.Kind().String())
			}
		case reflect.Int:
			if isAnyOfThese(ft.Kind(), reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64) {
				if !(v.(int) >= 0) {
					return false, BadFieldType(k, t.String(), ft.Kind().String())
				}
			} else if !isAnyOfThese(ft.Kind(), reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64) {
				return false, BadFieldType(k, t.String(), ft.Kind().String())
			}
		case reflect.Float32, reflect.Float64:
			if !isAnyOfThese(ft.Kind(), reflect.Float32, reflect.Float64) {
				return false, BadFieldType(k, t.String(), ft.Kind().String())
			}
		case reflect.Array, reflect.Slice:
			if reflect.TypeOf(v).AssignableTo(f.Type) {
				return false, BadFieldType(k, ft.String(), reflect.TypeOf(v).String())
			}
		case reflect.Map:
			if ft.Kind() != reflect.Struct {
				return false, BadFieldType(k, t.String(), ft.Kind().String())
			}
			return validate(v.(map[string]any), ft)
		default:
			return false, BadFieldType(k, "unsupported type: "+t.String(), ft.Kind().String())
		}
	}
	for k, required := range expectedNotFoundKeys {
		if required == true {
			return false, ExpectedButNotFound(k)
		}
	}
	return true, nil
}

// GetRequestHeaderParam gets the value of the param from the header. It may return error
// if it not found AND required.
func GetRequestHeaderParam(c *Context, param Parameter) (string, error) {
	value := c.GetRequestHeader(param.Name)
	return handleParam(value, param)
}

// getQueryParam gets the value of the param from the query. It may return error
// if it not found AND required.
func getQueryParam(c *Context, param Parameter) (string, error) {
	//FIXME: only the first letter should be lowered.
	value := c.GetQueryParam(param.Name)
	return handleParam(value, param)
}

// GetRequestHeaderParam gets the value of the param from the cookie header.
// It may return an error if it not found AND required.
func getCookieParam(c *Context, param Parameter) (string, error) {
	value := c.GetCookie(param.Name)
	return handleParam(value, param)
}

func getPathParam(index int, param Parameter, matches []string) (string, error) {
	if len(matches) > 1+index {
		m := matches[1+index]
		return handleParam(m, param)
	} else {
		return "", fmt.Errorf("required path param %s not provided", param.Name)
	}
}

// getBodyParam gets the value of the param from the body.
// It will return an error if it is not found AND required.
func getBodyParam(c *Context, param Parameter) (string, error) {
	// Read the body content
	body, err := c.GetBody()
	if err != nil {
		return "", fmt.Errorf("an error occurred while reading the body: %s", err.Error())
	}
	return handleParam(string(body), param)
}

func populateField(value string, field reflect.Value) error {
	fieldType := field.Type()
	switch fieldType.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		valuei, err := strconv.Atoi(value)
		if err != nil {
			return FieldTypeError(value, fieldType.Kind().String())
		}
		switch fieldType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(int64(valuei))
		default:
			if valuei < 0 {
				return FieldTypeError(value, fieldType.Kind().String())
			}
			valueui := uint64(valuei)
			field.SetUint(valueui)
		}
	case reflect.Float32, reflect.Float64:
		valuef, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return FieldTypeError(value, "float32")
		}
		field.SetFloat(valuef)
	case reflect.Bool:
		valueb, err := strconv.ParseBool(value)
		if err != nil {
			return FieldTypeError(value, "boolean")
		}
		field.SetBool(valueb)
	case reflect.Struct:
		var m map[string]any

		err := json.Unmarshal([]byte(value), &m)
		if err != nil {
			return InvalidJSON(value)
		}

		ok, err := validate(m, fieldType)
		if !ok {
			return err
		}

		newField := reflect.New(fieldType)
		err = json.Unmarshal([]byte(value), newField.Interface())
		if err != nil {
			return FieldTypeError(value, fieldType.Name())
		}
		field.Set(newField.Elem())
	}

	return nil
}

func populateInputSchema(c *Context, s any, p []Parameter, matches []string) error {
	if len(p) == 0 { //no input schema
		return nil
	}
	sve := reflect.ValueOf(s).Elem() //will not panic because we can confirm
	pathparamsindex := 0             //pathparamsindex is the amount of path params already reviewed
	for i, pa := range p {
		var value string
		var err error

		switch pa.In {
		case "header":
			value, err = GetRequestHeaderParam(c, pa)
		case "path":
			value, err = getPathParam(pathparamsindex, pa, matches)
		case "query":
			value, err = getQueryParam(c, pa)
		case "cookie":
			value, err = getCookieParam(c, pa)
		case "body":
			value, err = getBodyParam(c, pa)
		}
		if err != nil {
			return err
		}
		field := sve.Field(i) //has to be there because handleInputSchema
		err = populateField(value, field)
		if err != nil {
			return err
		}
	}
	return nil
}

// FIXME: type info lowercase
type typeInfo struct {
	_type string
	info  Schema
}

func newTypeInfo(_type string, s Schema) typeInfo {
	return typeInfo{
		_type: _type,
		info:  s,
	}
}

var supportedTypes = map[string]typeInfo{
	"string": newTypeInfo("string", Schema{
		Format:  "string",
		Example: "string",
	}),
	"int": newTypeInfo("integer", Schema{
		Format:  "int",
		Example: "0",
	}),
	"int8": newTypeInfo("number", Schema{
		// https://spec.openapis.org/registry/format/int8
		Format:  "int8",
		Example: "0",
	}),
	"int16": newTypeInfo("number", Schema{
		// https://spec.openapis.org/registry/format/int16
		Format:  "int16",
		Example: "0",
	}),
	"int32": newTypeInfo("number", Schema{
		// https://spec.openapis.org/registry/format/int32
		Format:  "int32",
		Example: "0",
	}),
	"int64": newTypeInfo("number", Schema{
		// https://spec.openapis.org/registry/format/int64
		Format:  "int64",
		Example: "0",
	}),
	"uint": newTypeInfo("integer", Schema{
		Format:  "int",
		Minimum: "0",
		Example: "0",
	}),
	"uint8": newTypeInfo("integer", Schema{
		Format:  "int8",
		Example: "0",
		Minimum: "0",
	}),
	"uint16": newTypeInfo("integer", Schema{
		Format:  "int16",
		Example: "0",
		Minimum: "0",
	}),
	"uint32": newTypeInfo("integer", Schema{
		Format:  "int32",
		Example: "0",
		Minimum: "0",
	}),
	"uint64": newTypeInfo("integer", Schema{
		Format:  "int64",
		Example: "0",
		Minimum: "0",
	}),
	"float32": newTypeInfo("number", Schema{
		Format:  "float",
		Example: "0.0",
	}),
	"float64": newTypeInfo("number", Schema{
		Format:  "double",
		Example: "0.0",
		Minimum: "0",
	}),
	"bool": newTypeInfo("boolean", Schema{
		Format:  "bool",
		Example: false,
	}),
}

func newDefinition(route *Route, schema any) Schema {
	newSchema := new(Schema)
	st := reflect.TypeOf(schema)
	sv := reflect.ValueOf(schema)

	if st.Kind() == reflect.Pointer {
		st = st.Elem()
		sv = sv.Elem()
	}

	// FIXME: refactor this it could look better
	if st.Kind() != reflect.Struct && st.Kind() != reflect.Slice && st.Kind() != reflect.Map && st.Kind() != reflect.Array && st.Kind() != reflect.Pointer {
		ts, ok := supportedTypes[st.String()]
		if !ok {
			panic("Unsupported type " + st.String() + ".")
		}
		newSchema.Type = ts.info.Type
		newSchema.Format = ts.info.Format
		newSchema.Minimum = ts.info.Minimum
		newSchema.Example = ts.info.Example
		return *newSchema
	}

	if st.Kind() == reflect.Map {
		if st.Key().Kind() != reflect.String {
			panic("map key type must always be string.")
		}
		nd := newDefinition(route, reflect.Zero(st.Elem()).Interface())
		newSchema.AdditionalProperties = &nd
		return *newSchema
	}
	if st.Kind() == reflect.Array || st.Kind() == reflect.Slice {
		newSchema.Type = "array"
		nd := newDefinition(route, reflect.Zero(st.Elem()).Interface())
		newSchema.Items = &nd
		return *newSchema
	}
	if st.Kind() != reflect.Struct {
		panic("type on field must either be string, int, bool, struct, slice, map, or array.")
	}
	// last remaining kind- reflect.Struct
	newDef := Schema{}
	newDef.Properties = make(map[string]*Schema)
	for i := range st.NumField() {
		newDef.Type = "object"
		field := st.Field(i)
		nd := newDefinition(route, sv.Field(i).Interface())

		fieldName := field.Name
		fieldNameSplit := strings.Split(field.Tag.Get("json"), ",")
		if len(fieldName) > 0 {
			fieldName = fieldNameSplit[0]
		}

		fieldRequired := field.Tag.Get("required")
		b, err := resolveBool(fieldRequired, true)
		if err != nil {
			panic(err)
		}
		if b == true {
			newDef.Required = append(newDef.Required, fieldName)
		}
		newDef.Properties[fieldName] = &nd
	}
	Definitions[st.Name()] = &newDef
	newSchema.Ref = "#/definitions/" + st.Name()
	return *newSchema
}
