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
		case http.MethodPut, http.MethodPatch, http.MethodDelete:
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
		return "text/plain" // default content type
	}
	return ct
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"error": message})
	w.WriteHeader(statusCode)
}

func handleJSONResponse(w http.ResponseWriter, req *http.Request, res JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res.Content); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func handleHTMLResponse(w http.ResponseWriter, req *http.Request, res HTMLResponse) {
	statusCode := resolveStatusCode(res.StatusCode, req.Method, res.Content)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	fmt.Fprint(w, res.Content)
}

func handleFileResponse(w http.ResponseWriter, req *http.Request, res FileResponse) {
	fileNameSplit := strings.Split(res.FileName, ".")
	suffix := fileNameSplit[len(fileNameSplit)-1]
	contentType := contentTypeFromFileSuffix(suffix)
	w.Header().Set("Content-Type", contentType)

	file, err := os.ReadFile(res.FileName)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Error retrieving file: "+err.Error())
		return
	}
	statusCode := resolveStatusCode(res.StatusCode, req.Method, string(file))
	w.Write(file)
	w.WriteHeader(statusCode)
}

func handleStreamingResponse(w http.ResponseWriter, res StreamingResponse) {
	// TODO: there should be a streaming struct to share data.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	stream := make(chan string)
	go func() {
		defer close(stream)
		res.StreamHandler(&stream)
	}()

	for value := range stream {
		fmt.Fprintf(w, "data: %s\n\n", value)
		w.(http.Flusher).Flush()
	}
}

func handleGenericResponse(w http.ResponseWriter, req *http.Request, res GenericResponse) {
	statusCode := resolveStatusCode(res.StatusCode, req.Method, res.Content)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	fmt.Fprint(w, res.Content)
}
