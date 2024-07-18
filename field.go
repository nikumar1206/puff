package puff

import (
	"fmt"
	"reflect"
)

type param struct {
	// modeled after https://swagger.io/specification/#:~:text=to%20the%20API.-,Fixed%20Fields,-Field%20Name
	Name        string
	Type        string // string, integer
	In          string //query, path, header, cookie
	Description string
	Required    bool
	Deprecated  bool
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

func handleInputSchema(s any) error { // should this return an error or should it panic?
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
			return fmt.Errorf("specified kind on field %s in struct tag must be header, path, query, body, or formData", svetf.Name)
		}

		//param.Description
		description := svetf.Tag.Get("description")

		//param.Required
		specified_required := svetf.Tag.Get("required")
		specified_deprecated := svetf.Tag.Get("deprecated")

		required, err := boolFromSpecified(specified_required, true)
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
	return nil
}
