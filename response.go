package puff

import "fmt"

type JSONResponse struct {
	StatusCode int
	Content    map[string]interface{}
}

func (j *JSONResponse) ResponseError(err error) string {
	return fmt.Sprintf("{\"error\": \"JSON Response Failed: %s\"}", err.Error())
}

type HTMLResponse struct { // the difference between this and Response is that the content type is text/html
	StatusCode int
	Content    string
}

type FileResponse struct {
	StatusCode      int
	FileName        string
	FileContentType string
}

func (f *FileResponse) Handler() func(Request) interface{} {
	return func(p Request) interface{} { return *f }
}

type StreamingResponse struct {
	StatusCode    int
	StreamHandler *func(*chan string)
}

type Response struct { // while this has a content-type of text/plain
	StatusCode  int
	Content     string
	ContentType string
}
