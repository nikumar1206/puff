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

func writeErrorResponse(ctx *Context, statusCode int, message string) {
	ctx.SetHeader("Content-Type", "application/json")
	json.NewEncoder(ctx.ResponseWriter).Encode(map[string]string{"error": message})
	ctx.SetStatusCode(statusCode)
}

func handleJSONResponse(ctx *Context, res JSONResponse) {
	ctx.SetHeader("Content-Type", "application/json")
	ctx.SetStatusCode(res.StatusCode)
	if err := json.NewEncoder(ctx.ResponseWriter).Encode(res.Content); err != nil {
		writeErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
}

func handleHTMLResponse(ctx *Context, res HTMLResponse) {
	statusCode := resolveStatusCode(res.StatusCode, ctx.Request.Method, res.Content)
	ctx.SetHeader("Content-Type", "text/html")
	ctx.SetStatusCode(statusCode)
	fmt.Fprint(ctx.ResponseWriter, res.Content)
}

func handleFileResponse(ctx *Context, res FileResponse) {
	fileNameSplit := strings.Split(res.FileName, ".")
	suffix := fileNameSplit[len(fileNameSplit)-1]
	contentType := contentTypeFromFileSuffix(suffix)
	ctx.SetHeader("Content-Type", contentType)

	file, err := os.ReadFile(res.FileName)
	if err != nil {
		writeErrorResponse(ctx, http.StatusInternalServerError, "Error retrieving file: "+err.Error())
		return
	}
	statusCode := resolveStatusCode(res.StatusCode, ctx.Request.Method, string(file))
	ctx.ResponseWriter.Write(file)
	ctx.SetStatusCode(statusCode)
}

func handleStreamingResponse(ctx *Context, res StreamingResponse) {
	// TODO: there should be a streaming struct to share data.
	ctx.SetHeader("Content-Type", "text/event-stream")
	ctx.SetHeader("Cache-Control", "no-cache")
	ctx.SetHeader("Connection", "keep-alive")

	stream := make(chan string)
	go func() {
		defer close(stream)
		res.StreamHandler(&stream)
	}()

	for value := range stream {
		fmt.Fprintf(ctx.ResponseWriter, "data: %s\n\n", value)
		ctx.ResponseWriter.(http.Flusher).Flush()
	}
}

func handleGenericResponse(ctx *Context, res GenericResponse) {
	statusCode := resolveStatusCode(res.StatusCode, ctx.Request.Method, res.Content)
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.SetStatusCode(statusCode)
	fmt.Fprint(ctx.ResponseWriter, res.Content)
}
