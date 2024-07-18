package puff

import (
	"fmt"
	"regexp"
)

type Route struct {
	fullPath string
	regexp   *regexp.Regexp
	params   []param

	Description string
	WebSocket   bool
	Protocol    string
	Pattern     string // TODO: remove this. un-necessary
	Path        string
	Handler     func(*Context)
	Fields      any
	// should probably have responses (200 OK followed by json, 400 Invalid Paramaters, etc...)
	// Responses []map[int]Response -> responses likely will look something like this
}

func (r *Route) String() string {
	return fmt.Sprintf("Protocol: %s\nPattern: %s\nPath: %s\n", r.Protocol, r.Pattern, r.Path)
}

func (r *Route) GetFullPath() string {
	return r.fullPath
}
