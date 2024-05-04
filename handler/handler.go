package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nikumar1206/puff/request"
	"github.com/nikumar1206/puff/response"
	"github.com/nikumar1206/puff/route"
)

func resolveStatusCode(sc int, method string) int {
	if sc == 0 {
		switch method {
		case http.MethodGet:
			return 200
		case http.MethodPost:
			return 201
		case http.MethodPut:
			return 204
		case http.MethodDelete:
			return 200
		default:
			return 200 // Default to 200 for unknown methods
		}
	}
	return sc
}

func resolveContentType(ct string) string {
	if ct == "" {
		return "text/plain"
	}
	return ct
}
func Handler(w http.ResponseWriter, req *http.Request, route *route.Route) {
	requestDetails := request.Request{}

	res := route.Handler(
		requestDetails,
	) // FIX ME: we should give the user handle function a request body as well
	var (
		contentType string
		content     string
		statusCode  int
	)
	switch r := res.(type) {
	case response.JSONResponse:
		statusCode = resolveStatusCode(r.StatusCode, req.Method)
		contentType = "application/json"
		contentBytes, err := json.Marshal(r.Content)
		if err != nil {
			content = r.ResponseError(err)
			http.Error(w, content, 500)
		}
		content = string(contentBytes)
	case response.HTMLResponse:
		statusCode = resolveStatusCode(r.StatusCode, req.Method)
		contentType = "text/html"
		content = r.Content
	case response.Response:
		statusCode = resolveStatusCode(r.StatusCode, req.Method)
		contentType = "text/plain"
		content = r.Content
	}

	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", contentType)
	fmt.Fprint(w, content)
}
