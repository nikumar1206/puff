package puff

import "fmt"

type Route struct {
	RouterName  string
	Protocol    string
	Pattern     string
	Path        string
	Description string
	Parameters  interface{}
	Handler     func(Request) interface{}

	// should probably have responses (200 OK followed by json, 400 Invalid Paramaters, etc...)
	// Responses []map[int]Response -> responses likely will look something like this
}

func (r *Route) String() string {
	return fmt.Sprintf("RouterName: %s\nProtocol: %s\nPattern: %s\nPath: %s\nDescription: %s\n",
		r.RouterName, r.Protocol, r.Pattern, r.Path, r.Description)
}
