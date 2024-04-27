package response

import "fmt"

// type baseResponse interface {
// }

type JSONResponse struct {
	StatusCode int
	Content    map[string]interface{}
}

func (j *JSONResponse) ResponseError(err string) (string, int) {
	return fmt.Sprintf("{\"error\": \"JSON Response Failed: %s\"}", err), 500
}

type HTMLResponse struct { //the difference between this and Response is that the content type is text/html
	StatusCode int
	Content    string
}
type Response struct { //while this has a content-type of text/plain
	StatusCode int
	Content    string
}
