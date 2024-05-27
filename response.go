package puff

type Response interface {
	GetStatusCode() int
	GetContent() any
}

type JSONResponse struct {
	StatusCode int
	Content    map[string]any
}

func (j JSONResponse) GetStatusCode() int {
	return j.StatusCode
}

func (j JSONResponse) GetContent() any {
	return j.Content
}

type HTMLResponse struct { // the difference between this and Response is that the content type is text/html
	StatusCode int
	Content    string
}

func (h HTMLResponse) GetStatusCode() int {
	return h.StatusCode
}

func (h HTMLResponse) GetContent() any {
	return h.Content
}

type GenericResponse struct { // while this has a content-type of text/plain
	StatusCode  int
	Content     string
	ContentType string
}

func (g GenericResponse) GetStatusCode() int {
	return g.StatusCode
}

func (g GenericResponse) GetContent() any {
	return g.Content
}
