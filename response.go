package puff

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
)

// Response is an interface that all response types should implement.
type Response interface {
	GetStatusCode() int
	GetContentType() string
	WriteContent(http.ResponseWriter, *http.Request) error
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

func (j JSONResponse) GetContentType() string {
	return "application/json"
}

// GetContent returns the content of the JSON response.
func (j JSONResponse) WriteContent(w http.ResponseWriter, r *http.Request) error {
	err := json.NewEncoder(w).Encode(j.Content)
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
	return h.StatusCode
}

func (h HTMLResponse) GetContentType() string {
	return "text/html"
}

// GetContent returns the content of the HTML response.
func (h HTMLResponse) WriteContent(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, h.Content)
	return nil
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

func (f FileResponse) GetContentType() string {
	return resolveContentType(f.ContentType, contentTypeFromFileName(f.FileName))
}

// GetContent returns the file content.
func (f FileResponse) WriteContent(w http.ResponseWriter, r *http.Request) error {
	file, err := os.ReadFile(f.FileName)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Error retrieving file: "+err.Error())
		return fmt.Errorf("Error retrieving file %s during FileResponse: %s", f.FileName, err.Error())
	}

	w.Write(file)
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
	StatusCode    int
	StreamHandler func(*chan string)
}

// GetStatusCode returns the status code of the streaming response.
func (s StreamingResponse) GetStatusCode() int {
	return s.StatusCode
}

func (s StreamingResponse) GetContentType() string {
	return "text/event-stream"
}

// GetContent returns the content of the streaming response.
func (s StreamingResponse) WriteContent(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	stream := make(chan string)
	go func() {
		defer close(stream)
		s.StreamHandler(&stream)
	}()

	for value := range stream {
		fmt.Fprintf(w, "data: %s\n\n", value)
		w.(http.Flusher).Flush()
	}
	return nil
}

func (s StreamingResponse) Handler() func(*Context) {
	return func(c *Context) {
		c.SendResponse(s)
	}
}

type WebSocketResponse struct {
}

func secWebSocketAcceptKey(key string) string {
	guid := "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.New()
	hash.Write([]byte(key + guid))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

type WebSocket struct {
	Context   *Context
	Conn      net.Conn
	Channel   chan string
	OnMessage func(WebSocket, string) //optional: events
	OnClose   func(WebSocket)         //optional: events
}

func read_conn(ws WebSocket) {
	for {
		buf := make([]byte, 1024)
		n, err := ws.Conn.Read(buf)
		if err != nil {
			break
		}
		ws.OnMessage(ws, string(buf[:n]))
	}
	ws.Close()
}

func (ws WebSocket) Send(message string) error {
	_, err := ws.Conn.Write([]byte(message))
	return err
}

func (ws WebSocket) Close() {
	ws.OnClose(ws)
	close(ws.Channel)
	ws.Conn.Close() //do not care about erros
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

func (g GenericResponse) GetContentType() string {
	return resolveContentType(g.ContentType, "text/plain")
}

// GetContent returns the content of the generic response.
func (g GenericResponse) WriteContent(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, g.Content)
	return nil
}
