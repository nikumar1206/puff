package puff

// Response is an interface that all response types should implement.
type Response interface {
	GetStatusCode() int
	GetContent() any
}

// JSONResponse represents a response with JSON content.
type JSONResponse struct {
	StatusCode int
	Content    map[string]any
}

// GetStatusCode returns the status code of the JSON response.
func (j JSONResponse) GetStatusCode() int {
	return j.StatusCode
}

// GetContent returns the content of the JSON response.
func (j JSONResponse) GetContent() any {
	return j.Content
}

// HTMLResponse represents a response with HTML content.
type HTMLResponse struct {
	StatusCode int
	Content    string
}

// GetStatusCode returns the status code of the HTML response.
func (h HTMLResponse) GetStatusCode() int {
	return h.StatusCode
}

// GetContent returns the content of the HTML response.
func (h HTMLResponse) GetContent() any {
	return h.Content
}

// FileResponse represents a response that sends a file.
type FileResponse struct {
	StatusCode  int
	FileName    string
	FileContent []byte
	ContentType string
}

// GetStatusCode returns the status code of the file response.
func (f FileResponse) GetStatusCode() int {
	return f.StatusCode
}

// GetContent returns the file content.
func (f FileResponse) GetContent() any {
	return f.FileContent
}

// Handler returns a handler function for serving the file response.
func (f *FileResponse) Handler() func(Request) Response {
	return func(p Request) Response {
		return *f
	}
}

// StreamingResponse represents a response that streams content.
type StreamingResponse struct {
	StatusCode    int
	StreamHandler func(*chan string)
	content       string
}

// GetStatusCode returns the status code of the streaming response.
func (s StreamingResponse) GetStatusCode() int {
	return s.StatusCode
}

// GetContent returns the content of the streaming response.
func (s StreamingResponse) GetContent() any {
	return s.content
}

// GenericResponse represents a response with plain text content.
type GenericResponse struct {
	StatusCode  int
	Content     string
	ContentType string
}

// GetStatusCode returns the status code of the generic response.
func (g GenericResponse) GetStatusCode() int {
	return g.StatusCode
}

// GetContent returns the content of the generic response.
func (g GenericResponse) GetContent() any {
	return g.Content
}
