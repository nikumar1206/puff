package puff

type Route struct {
	RouterName  string
	Protocol    string
	Pattern     string
	Path        string
	Description string
	Handler     func(Request) interface{}
	// should probably have responses (200 OK followed by json, 400 Invalid Paramaters, etc...)
	// Responses []map[int]Response -> responses likely will look something like this
}
