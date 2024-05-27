package puff

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		content, err := json.Marshal(map[string]string{"message": err.Error()})

		if err != nil {
			panic(err)
		}

		if err != nil {
			http.Error(w, string(content), 500)
		}
		return
	case HTMLResponse:
		statusCode = resolveStatusCode(r.StatusCode, req.Method)
		contentType = "text/html"
		content = r.Content
	case GenericResponse:
		statusCode = resolveStatusCode(r.StatusCode, req.Method)
		contentType = "text/plain"
		content = r.Content
	default:
		http.Error(w, "The response type provided to handle this request is invalid.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", contentType)
	fmt.Fprint(w, content)
}
