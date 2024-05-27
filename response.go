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

// returns html content with 'Content-Type' header set to text/html
type HTMLResponse struct {
	StatusCode int
	Content    string
}

// returns a
// FIXME: feels wrong to send a file as string instead of []byte
type FileResponse struct {
	StatusCode      int
	FileName        string
	FileContentType string
}

func (f FileResponse) GetStatusCode() int {
	return f.StatusCode
}

func (f FileResponse) GetContent() any {
	return f.FileContentType
}

func (f *FileResponse) Handler() func(Request) Response {
	return func(p Request) Response { return *f }
}

// to be used when sending server side events
type StreamingResponse struct {
	StatusCode    int
	StreamHandler func(*chan string)
	content       string
}

func (s StreamingResponse) GetStatusCode() int {
	return s.StatusCode
}

func (s StreamingResponse) GetContent() any {
	return s.content
}

func (h HTMLResponse) GetStatusCode() int {
	return h.StatusCode
}

func (h HTMLResponse) GetContent() any {
	return h.Content
}

// default response structure with 'Content-Type' header set to text/plain
type GenericResponse struct {
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
