package puff

import (
	"fmt"
	"regexp"
)

type Route struct {
	WebSocket bool
	Protocol  string
	Path      string
	Handler   func(*Context)
	fullPath  string
	regexp    *regexp.Regexp
	Fields    Field
	// should probably have responses (200 OK followed by json, 400 Invalid Paramaters, etc...)
	// Responses []map[int]Response -> responses likely will look something like this
}

func (r *Route) String() string {
	return fmt.Sprintf("Protocol: %s\nPath: %s\n", r.Protocol, r.Path)
}

func (r *Route) GetFullPath() string {
	return r.fullPath
}
