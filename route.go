package puff

import (
	"fmt"
	"maps"
	"reflect"
	"regexp"
	"strings"
)

type Route struct {
	fullPath string
	regexp   *regexp.Regexp
	params   []Parameter

	Description string
	WebSocket   bool
	Protocol    string
	Path        string
	Handler     func(*Context)
	Fields      any
	// Router points to the router the route belongs to. Will always be the closest router in the tree.
	Router *Router
	// should probably have responses (200 OK followed by json, 400 Invalid Paramaters, etc...)
	Responses map[int]Response
}

func (r *Route) String() string {
	return fmt.Sprintf("Protocol: %s\nPath: %s\n", r.Protocol, r.Path)
}

func (r *Route) GetFullPath() string {
	return r.fullPath
}

func (route *Route) getCompletePath() {
	var parts []string
	currentRouter := route.Router
	for currentRouter != nil {
		parts = append([]string{currentRouter.Prefix}, parts...)
		currentRouter = currentRouter.parent
	}

	parts = append(parts, route.Path)
	route.fullPath = strings.Join(parts, "")
}

func (route *Route) createRegexMatch() {
	escapedPath := strings.ReplaceAll(route.fullPath, "/", "\\/")
	regexPattern := regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(escapedPath, "([^/]+)")
	route.regexp = regexp.MustCompile("^" + regexPattern + "$")
}

func (route *Route) handleInputSchema() error { // should this return an error or should it panic?
	if route.Fields == nil {
		*&route.params = []Parameter{}
		return nil
	}
	sv := reflect.ValueOf(route.Fields) //
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
		newParam.Schema = newDefinition(route, sve.Field(i).Interface())

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

		required, err := resolveBool(specified_required, required_def)
		if err != nil {
			return err
		}
		deprecated, err := resolveBool(specified_deprecated, false)
		if err != nil {
			return err
		}

		//param.Schema.format
		format := svetf.Tag.Get("format")
		if format != "" {
			newParam.Schema.Format = format
		}

		newParam.Name = name
		newParam.In = specified_kind
		newParam.Description = description
		newParam.Required = required
		newParam.Deprecated = deprecated

		newParams = append(newParams, newParam)
	}
	route.params = newParams
	return nil
}

// GenerateResponses is responsible for generating the 'responses' attribute in the OpenAPI schema.
// Since responses can be specified at multiple levels, responses at the route level will be given the most specificity.
func (r *Route) GenerateResponses() {

	if r.Router.puff.DocsURL == "" {
		// if swagger documentation is off, we will not set responses
		return
	}
	responses := r.Responses
	if responses == nil {
		responses = make(map[int]Response)
	}
	currentRouter := r.Router

	for currentRouter != nil {
		// avoid over-writing the original responses for the routers
		clonedResponses := maps.Clone(currentRouter.Responses)
		fmt.Println("preclone", clonedResponses)
		fmt.Println("what is cloned", clonedResponses, responses)
		maps.Copy(clonedResponses, responses)
		currentRouter = currentRouter.parent
	}

	r.Responses = responses
}
