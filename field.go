package puff

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// FIXME: Allow for the names of params to be different than the names of the structfields.
// maybe we could even automatically lower it?

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

// boolFromSpecified resolves the specified bool (as a string type)
// and the default bool. It gives priority to the specified.
func boolFromSpecified(spec string, def bool) (bool, error) {
	var b bool
	switch spec {
	case "":
		b = def
	case "true":
		b = true
	case "false":
		b = false
	default:
		return false, fmt.Errorf("specified required on field must be either true or false")
	}
	return b, nil
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

// getHeaderParam gets the value of the param from the header. It may return error
// if it not found AND required.
func getHeaderParam(c *Context, param Parameter) (string, error) {
	value := c.GetHeader(param.Name)
	return handleParam(value, param)
}

// getQueryParam gets the value of the param from the query. It may return error
// if it not found AND required.
func getQueryParam(c *Context, param Parameter) (string, error) {
	//FIXME: only the first letter should be lowered.
	value := c.GetQueryParam(param.Name)
	return handleParam(value, param)
}

// getHeaderParam gets the value of the param from the cookie header.
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

func setStringField(sf reflect.Value, body []byte, param Parameter) error {
	content, err := handleParam(string(body), param)
	if err != nil {
		return err
	}
	sf.SetString(content)
	return nil
}
func setIntField(sf reflect.Value, body []byte, param Parameter) error {
	content, err := handleParam(string(body), param)
	if err != nil {
		return err
	}
	ci, err := strconv.Atoi(content)
	if err != nil {
		return err
	}
	sf.SetInt(int64(ci))
	return nil
}

// getBodyParam gets the value of the param from the body.
// It will return an error if it is not found AND required.
func getBodyParam(c *Context, param Parameter) (string, error) {
	// Read the body content
	defer c.Request.Body.Close()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return "", fmt.Errorf("an error occurred while reading the body: %s", err.Error())
	}
	return handleParam(string(body), param)
}

func populateField(value string, field reflect.Value) error {
	fieldType := field.Type()
	newField := reflect.New(fieldType)

	err := json.Unmarshal([]byte(value), newField.Interface())
	if err != nil {
		return err
	}
	field.Set(newField.Elem())
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
			value, err = getHeaderParam(c, pa)
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

var supportedTypes = map[string]string{
	"string": "string",
	"int":    "integer",
	// FIXME: int8, int16, int32, int64 are not the same thing as int
	"int8":  "integer",
	"int16": "integer",
	"int32": "integer",
	"int64": "integer",
	// FIXME: uint, uint8, uint16, uint32, uint64 are not the same thing as int
	"uint":   "integer",
	"uint8":  "integer",
	"uint16": "integer",
	"uint32": "integer",
	"uint64": "integer",
	// FIXME: float32 and float64 are not the same thing
	"float32": "number",
	"float64": "number",
	"bool":    "boolean",
}

func newDefinition(schema any) Schema {
	newSchema := new(Schema)
	st := reflect.TypeOf(schema)
	sv := reflect.ValueOf(schema)
	// FIXME: refactor this it could look better
	if st.Kind() != reflect.Struct && st.Kind() != reflect.Slice && st.Kind() != reflect.Map && st.Kind() != reflect.Array && st.Kind() != reflect.Pointer {
		// FIXME: st.String will return int even though the specification uses integer.
		ts, ok := supportedTypes[st.String()]
		if !ok {
			panic("Unsupported type " + st.String() + ".")
		}
		newSchema.Type = ts
		return *newSchema
	}
	if st.Kind() == reflect.Pointer {
		panic("Pointers are not supported.")
	}
	if st.Kind() == reflect.Map {
		if st.Key().Kind() != reflect.String {
			panic("Map key type must always be string.")
		}
		nd := newDefinition(reflect.Zero(st.Elem()).Interface())
		newSchema.AdditionalProperties = &nd
		return *newSchema
	}
	if st.Kind() == reflect.Array || st.Kind() == reflect.Slice {
		newSchema.Type = "array"
		nd := newDefinition(reflect.Zero(st.Elem()).Interface())
		newSchema.Items = &nd
		return *newSchema
	}
	if st.Kind() != reflect.Struct {
		panic("Type on field must either be string, int, bool, struct, slice, map, or array.")
	}
	// last remaining kind- reflect.Struct
	newDef := Schema{}
	newDef.Properties = make(map[string]*Schema)
	for i := range st.NumField() {
		newDef.Type = "object"
		field := st.Field(i)
		nd := newDefinition(sv.Field(i).Interface())

		fieldName := field.Name
		fieldNameSplit := strings.Split(field.Tag.Get("json"), ",")
		if len(fieldName) > 0 {
			fieldName = fieldNameSplit[0]
		}
		newDef.Properties[fieldName] = &nd
	}
	AddDefinition(st.Name(), newDef)
	newSchema.Ref = "#/definitions/" + st.Name()
	return *newSchema
}

func handleInputSchema(pa *[]Parameter, s any) error { // should this return an error or should it panic?
	if s == nil {
		*pa = []Parameter{}
		return nil
	}
	sv := reflect.ValueOf(s) //
	svk := sv.Kind()
	if svk != reflect.Ptr {
		return fmt.Errorf("fields must be POINTER to struct")
	}
	sve := sv.Elem()
	svet := sve.Type()
	if sve.Kind() != reflect.Struct {
		return fmt.Errorf("fields must be pointer to STRUCT")
	}

	newParams := []Parameter{}
	for i := range svet.NumField() {
		newParam := Parameter{}
		svetf := svet.Field(i)

		name := svetf.Tag.Get("name")
		if name == "" {
			name = svetf.Name
		}

		// param.Schema
		newParam.Schema = newDefinition(sve.Field(i).Interface())

		//param.In
		specified_kind := svetf.Tag.Get("kind") //ref: Parameters object/In
		if name == "Body" && specified_kind == "" {
			specified_kind = "body"
		}
		if !isValidKind(specified_kind) {
			return fmt.Errorf("specified kind on field %s in struct tag must be header, path, query, cookie, body, or formdata", svetf.Name)
		}

		//param.Description
		description := svetf.Tag.Get("description")

		//param.Required
		specified_required := svetf.Tag.Get("required")
		specified_deprecated := svetf.Tag.Get("deprecated")

		required_def := true
		if specified_kind == "cookie" { // cookies by default should never be required
			required_def = false
		}

		required, err := boolFromSpecified(specified_required, required_def)
		if err != nil {
			return err
		}
		deprecated, err := boolFromSpecified(specified_deprecated, false)
		if err != nil {
			return err
		}

		//param.Schema.format
		format := svetf.Tag.Get("format")
		newParam.Name = name
		newParam.In = specified_kind
		newParam.Schema.Format = format
		newParam.Description = description
		newParam.Required = required
		newParam.Deprecated = deprecated

		newParams = append(newParams, newParam)
	}
	*pa = newParams
	return nil
}
