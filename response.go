package puff

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Response is an interface that all response types should implement.
type Response interface {
	GetStatusCode() int
	GetContentType() string
	WriteContent(*Context) error
}

// JSONResponse represents a response with JSON content.
type JSONResponse struct {
	StatusCode int
	Content    any
}

// GetStatusCode returns the status code of the JSON response.
func (j JSONResponse) GetStatusCode() int {
	return resolveStatusCode(j.StatusCode, 200)
}

func (j JSONResponse) GetContentType() string {
	return "application/json"
}

// GetContent returns the content of the JSON response.
func (j JSONResponse) WriteContent(c *Context) error {
	err := json.NewEncoder(c.ResponseWriter).Encode(j.Content)
	if err != nil {
		return fmt.Errorf("Writing JSONResponse Content failed with: %s", err.Error())
	}
	return nil
}

// HTMLResponse represents a response with HTML content.
type HTMLResponse struct {
	StatusCode int
	Content    string
}

// GetStatusCode returns the status code of the HTML response.
func (h HTMLResponse) GetStatusCode() int {
	return resolveStatusCode(h.StatusCode, 200)
}

func (h HTMLResponse) GetContentType() string {
	return "text/html"
}

// GetContent returns the content of the HTML response.
func (h HTMLResponse) WriteContent(c *Context) error {
	fmt.Fprint(c.ResponseWriter, h.Content)
	return nil
}

// FileResponse represents a response that sends a file.
type FileResponse struct {
	StatusCode  int
	FilePath    string
	FileContent []byte
	ContentType string
}

// GetStatusCode returns the status code of the file response.
func (f FileResponse) GetStatusCode() int {
	return resolveStatusCode(f.StatusCode, 200)
}

func (f FileResponse) GetContentType() string {
	return resolveContentType(f.ContentType, contentTypeFromFileName(f.FilePath))
}

// GetContent returns the file content.
func (f FileResponse) WriteContent(c *Context) error {
	file, err := os.ReadFile(f.FilePath)
	if err != nil {
		writeErrorResponse(
			c.ResponseWriter,
			http.StatusInternalServerError,
			"Error retrieving file: "+err.Error(),
		)
		return fmt.Errorf(
			"error retrieving file %s during FileResponse: %s",
			f.FilePath,
			err.Error(),
		)
	}

	c.ResponseWriter.Write(file)
	return nil
}

// Handler returns a handler function for serving the file response.
func (f *FileResponse) Handler() func(*Context) {
	return func(c *Context) {
		c.SendResponse(f)
	}
}

// StreamingResponse represents a response that streams content.
type StreamingResponse struct {
	StatusCode int
	// StreamHandler is a function that takes in a pointer to a channel.
	// The channel should be written to with a ServerSideEvent to write
	// to the response. It should be closed once done writing.
	StreamHandler func(*chan ServerSideEvent)
}

type ServerSideEvent struct {
	Event string
	Data  string
	ID    string
	Retry int
}

// GetStatusCode returns the status code of the streaming response.
func (s StreamingResponse) GetStatusCode() int {
	return resolveStatusCode(s.StatusCode, 200)
}

func (s StreamingResponse) GetContentType() string {
	return "text/event-stream"
}

// GetContent returns the content of the streaming response.
func (s StreamingResponse) WriteContent(c *Context) error {
	c.ResponseWriter.Header().Set("Cache-Control", "no-cache")
	c.ResponseWriter.Header().Set("Connection", "keep-alive")

	stream := make(chan ServerSideEvent)
	go func() {
		defer close(stream)
		s.StreamHandler(&stream)
	}()
	for value := range stream {
		fmt.Fprint(c.ResponseWriter, constructSSE(value))
		c.ResponseWriter.(http.Flusher).Flush()
	}
	return nil
}

func constructSSE(eventStruct ServerSideEvent) string {
	finalEvent := ""

	if eventStruct.ID != "" {
		finalEvent += fmt.Sprintf("id: %s\n", eventStruct.ID)
	}

	if eventStruct.Event != "" {
		finalEvent += fmt.Sprintf("event: %s\n", eventStruct.Event)
	}

	if eventStruct.Retry != 0 {
		finalEvent += fmt.Sprintf("retry: %d\n", eventStruct.Retry)
	}

	finalEvent += fmt.Sprintf("data: %s\n\n", eventStruct.Data)

	return finalEvent
}

func (s StreamingResponse) Handler() func(*Context) {
	return func(c *Context) {
		c.SendResponse(s)
	}
}

// RedirectResponse represents a response that sends a redirect to the client.
type RedirectResponse struct {
	// StatusCode provides the 3xx status code of the redirect response. Default: 308.
	StatusCode int
	// To provides the URL to redirect the client to.
	To string
}

// GetStatusCode returns the status code for the redirect response.
// If the status code is not provided, or not valid for a redirect, it will default to 308.
func (r RedirectResponse) GetStatusCode() int {
	if r.StatusCode == 0 || !(r.StatusCode >= 300 && r.StatusCode <= 308) {
		return 308
	}
	return r.StatusCode
}

// GetContentType returns the content type for the redirect response.
// It will always return an empty string since there is no body to
// describe in content typc.ResponseWriter.
func (r RedirectResponse) GetContentType() string {
	return "text/html; charset=utf-8"
}

// WriteContent writes the header Location to redirect the client to.
func (r RedirectResponse) WriteContent(c *Context) error {
	c.ResponseWriter.Header().Set("Location", r.To)
	fmt.Fprintf(c.ResponseWriter, `<!DOCTYPE HTML>
    <html lang='en-US'>
    <head>
        <meta charset='UTF-8'>
        <meta http-equiv='refresh' content='0; url=%s'>
        <script type='text/javascript'>
            window.location.href = '%s'
        </script>
        <title>Page Redirection</title>
    </head>
    <body>
        If you are not redirected automatically, follow this <a href='%s'>link to example</a>.
    </body>
    </html>`, r.To, r.To, r.To)
	return nil
}

// GenericResponse represents a response with plain text content.
type GenericResponse struct {
	StatusCode  int
	Content     string
	ContentType string
}

// GetStatusCode returns the status code of the generic response.
func (g GenericResponse) GetStatusCode() int {
	return resolveStatusCode(g.StatusCode, 200)
}

func (g GenericResponse) GetContentType() string {
	return resolveContentType(g.ContentType, "text/plain")
}

// GetContent returns the content of the generic response.
func (g GenericResponse) WriteContent(c *Context) error {
	fmt.Fprint(c.ResponseWriter, g.Content)
	return nil
}
