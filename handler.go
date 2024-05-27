package puff

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"
)

func resolveStatusCode(sc int, method string, content any) int {
	if content == "" {
		return http.StatusNoContent
	}
	if sc == 0 {
		switch method {
		case http.MethodGet:
			return http.StatusOK
		case http.MethodPost:
			return http.StatusCreated
		case http.MethodPut:
			return http.StatusOK
		case http.MethodPatch:
			return http.StatusOK
		case http.MethodDelete:
			return http.StatusOK
		default:
			return http.StatusOK
		}
	}
	return sc
}

func contentTypeFromFileSuffix(suffix string) string {
	ct := mime.TypeByExtension("." + suffix)
	if ct == "" {
		return "text/plain" //we dont know the content type from file suffix
	}
	return ct
}

func Handler(w http.ResponseWriter, req *http.Request, route *Route) {
	requestDetails := Request{}

	res := route.Handler(
		requestDetails,
	) // FIX ME: we should give the user handle function a request body as well

	var (
		contentType string
		content     string
		statusCode  int
	)
	switch r := res.(type) {
	case JSONResponse:
		contentType = "application/json"
		w.Header().Set("Content-Type", contentType)
		statusCode = resolveStatusCode(r.StatusCode, req.Method, r.GetContent())
		err := json.NewEncoder(w).Encode(r.Content)
		content, err := json.Marshal(map[string]string{"message": err.Error()})
		if err != nil {
			panic(err)
		}
		if err != nil {
			http.Error(w, string(content), 500)
		}
		w.WriteHeader(statusCode)
		return
	case HTMLResponse:
		statusCode = resolveStatusCode(r.StatusCode, req.Method)
		contentType = "text/html"
		content = r.Content
	case FileResponse:
		fileNameSplit := strings.Split(r.FileName, ".")
		suffix := fileNameSplit[len(fileNameSplit)-1]
		contentType = contentTypeFromFileSuffix(suffix)
		file, err := os.ReadFile(r.FileName)
		if err != nil {
			statusCode = 500
			content = "There was an error retrieving the file: " + err.Error()
		}
		statusCode = resolveStatusCode(r.StatusCode, req.Method)
		content = string(file)
	case StreamingResponse:
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		stream := make(chan string)
		go func() {
			defer close(stream)
			r.StreamHandler(&stream)
		}()
		for value := range stream {
			fmt.Fprintf(w, "data: %s\n\n", value)
			w.(http.Flusher).Flush()
		}
		return
	case GenericResponse:
		statusCode = resolveStatusCode(r.StatusCode, req.Method)
		content = r.Content
	default:
		http.Error(w, "The response type provided to handle this request is invalid.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", contentType)
	fmt.Fprint(w, content)
}
