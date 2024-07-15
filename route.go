package puff

import (
	"fmt"
	"regexp"
)

type Route struct {
	Protocol string
	Pattern  string // TODO: remove this. un-necessary
	Path     string
	Handler  func(*Context)
	fullPath string
	regexp   *regexp.Regexp
	Fields   Field
	// should probably have responses (200 OK followed by json, 400 Invalid Paramaters, etc...)
	// Responses []map[int]Response -> responses likely will look something like this
}

func (r *Route) String() string {
	return fmt.Sprintf("Protocol: %s\nPattern: %s\nPath: %s\nFullPath: %s\n",

		r.Protocol, r.Pattern, r.Path, r.fullPath)
}

func (r *Route) GetFullPath() string {
	return r.fullPath
}
