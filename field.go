package puff

import (
	"fmt"
	"reflect"
	"strconv"
)

type param struct {
	// modeled after https://swagger.io/specification/#:~:text=to%20the%20API.-,Fixed%20Fields,-Field%20Name
	Name        string `json:"name"`
	Type        string `json:"type"` // string, integer
	In          string `json:"in"`   //query, path, header, cookie
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Deprecated  bool   `json:"deprecated"`
}

// type HelloWorld struct {
// 	Name int `kind:"QueryParam" description:"Specify a name to issue a greeting to."`
// }
//
// func main() {
// 	var router Router
// 	HelloWorldInput := new(HelloWorld)
// 	router.Get("/hello-world", HelloWorldInput, func(c *Context) {
// 		c.SendResponse(GenericResponse{
// 			Content: fmt.Sprintf("Hello, %s!", HelloWorldInput.Name),
// 		})
// 	})
// }

func isValidType(_type string) bool {
	// https://swagger.io/specification/#:~:text=openapi.yaml.-,Data%20Types,-Data%20types%20in
	return _type == "string" || _type == "int"
}

func isValidKind(specified_kind string) bool {
	// https://swagger.io/specification/#:~:text=the%20in%20property.-,in,query%22%2C%20%22header%22%2C%20%22path%22%20or%20%22cookie%22.,-description
	// although the spec only specifies these types, the other types are written in the schema petstore example
	// meaning these types are supported
	// https://petstore.swagger.io/v2/swagger.json

	return specified_kind == "header" ||
		specified_kind == "path" ||
		specified_kind == "query" ||
		specified_kind == "cookie" ||
		specified_kind == "body" ||
		specified_kind == "formData"
}

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
func handleParam(value string, param param) (string, error) {
	ok := !(value == "")
	if !ok && param.Required {
		return "", fmt.Errorf("Required %s param %s not provided.", param.In, param.Name)
	}
	return value, nil
}

// getHeaderParam gets the value of the param from the header. It may return error
// if it not found AND required.
func getHeaderParam(c *Context, param param) (string, error) {
	value := c.GetHeader(param.Name)
	return handleParam(value, param)
}

// getQueryParam gets the value of the param from the query. It may return error
// if it not found AND required.
func getQueryParam(c *Context, param param) (string, error) {
	//FIXME: only the first letter should be lowered.
	value := c.GetQueryParam(param.Name)
	return handleParam(value, param)
}

// getHeaderParam gets the value of the param from the cookie header.
// It may return an error if it not found AND required.
func getCookieParam(c *Context, param param) (string, error) {
	value := c.GetCookie(param.Name)
	return handleParam(value, param)
}

func getPathParam(c *Context, index int, param param, matches []string) (string, error) {
	if len(matches) > index {
		m := matches[index]
		return handleParam(m, param)
	} else {
		fmt.Println(index, param, matches)
		return "", fmt.Errorf("required path param %s not provided", param.Name)
	}
}

// getBodyParam gets the value of the param from the body.
// It will return an error if it is not found AND required.

func populateInputSchema(c *Context, s any, p []param, matches []string) error {
	if len(p) == 0 { //no input schema
		return nil
	}
	sve := reflect.ValueOf(s).Elem()
	pathparamsindex := 0
	for _, pa := range p {
		var value string
		var err error
		switch pa.In {
		case "header":
			value, err = getHeaderParam(c, pa)
		case "path":
			value, err = getPathParam(c, pathparamsindex, pa, matches)
			// continue // FIXME: how do i get path?
		case "query":
			value, err = getQueryParam(c, pa)
		case "cookie":
			value, err = getCookieParam(c, pa)
		}
		if err != nil {
			return err
		}

		field := sve.FieldByName(pa.Name) //has to be there because handleInputSchema
		if pa.Type == "int" {
			valuei, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("Expected type integer on param %s.", pa.Name)
			}
			field.SetInt(int64(valuei))
			continue
		}
		field.SetString(value)
	}
	return nil
}

func handleInputSchema(pa *[]param, s any) error { // should this return an error or should it panic?
	if s == nil {
		*pa = []param{}
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
		return fmt.Errorf("fields must be pointer to struct")
	}

	newParams := []param{}
	for i := range svet.NumField() {
		newParam := param{}
		svetf := svet.Field(i)

		// param.Type
		_type := svetf.Type.String()
		if !isValidType(_type) {
			return fmt.Errorf("type on field %s must be string or int", svetf.Name)
		}

		//param.In
		specified_kind := svetf.Tag.Get("kind") //ref: Parameters object/In
		if !isValidKind(specified_kind) {
			return fmt.Errorf("specified kind on field %s in struct tag must be header, path, query, cookie, body, or formData", svetf.Name)
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

		newParam.Name = svetf.Name
		newParam.In = specified_kind
		newParam.Type = _type
		newParam.Description = description
		newParam.Required = required
		newParam.Deprecated = deprecated

		newParams = append(newParams, newParam)
	}
	*pa = newParams
	return nil
}
