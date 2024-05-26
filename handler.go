package puff

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"
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

func resolveContentType(ct string, mct string) string {
	if ct == "" {
		return mct
	}
	return ct
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
		statusCode = resolveStatusCode(r.StatusCode, req.Method)
		w.Header().Add("Content-Type", contentType)
		w.WriteHeader(statusCode)
		err := json.NewEncoder(w).Encode(r.Content)
		if err != nil {
			content = r.ResponseError(err)
			http.Error(w, content, 500)
		}
		return
	case HTMLResponse:
		statusCode = resolveStatusCode(r.StatusCode, req.Method)
		contentType = "text/html"
		content = r.Content
	case FileResponse:
		if r.FileContentType == "" {
			fileNameSplit := strings.Split(r.FileName, ".")
			suffix := fileNameSplit[len(fileNameSplit)-1]
			contentType = contentTypeFromFileSuffix(suffix)
		} else {
			contentType = r.FileContentType
		}
		file, err := os.ReadFile(r.FileName)
		if err != nil {
			statusCode = 500
			content = "There was an error retrieving the file: " + err.Error()
		}
		statusCode = resolveStatusCode(r.StatusCode, req.Method)
		content = string(file)
	case Response:
		statusCode = resolveStatusCode(r.StatusCode, req.Method)
		contentType = resolveContentType(r.ContentType, "text/plain")
		content = r.Content
	default:
		http.Error(w, "The response type provided to handle this request is invalid.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", contentType)
	fmt.Fprint(w, content)
}
