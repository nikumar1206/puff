package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	response "puff/response"
)

func resolveStatusCode(sc int) int {
	if sc == 0 {
		return 200
	}
	return sc
}

func Handler(w http.ResponseWriter, req *http.Request, handlerFunc func() interface{}) {
	// FIX ME: middleware comes here
	res := handlerFunc() // FIX ME: we should give the user handle function a request body as well
	var (
		contentType string
		content     string
		statusCode  int
	)
	switch r := res.(type) {
	case response.JSONResponse:
		statusCode = resolveStatusCode(r.StatusCode)
		contentType = "application/json"
		contentBytes, err := json.Marshal(r.Content)
		if err != nil {
			content, statusCode = r.ResponseError(err.Error())
			break
		}
		content = string(contentBytes)
	case response.HTMLResponse:
		statusCode = resolveStatusCode(r.StatusCode)
		contentType = "text/html"
		content = r.Content
	case response.Response:
		statusCode = resolveStatusCode(r.StatusCode)
		contentType = "text/plain"
		content = r.Content
	}

	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", contentType)
	fmt.Fprint(w, content)
}
