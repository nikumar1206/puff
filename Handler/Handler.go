package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	response "puff/Response"
	router "puff/Router"
)

func resolveStatusCode(sc int) int {
	if sc == 0 {
		return 200
	}
	return sc
}

func Handler(w http.ResponseWriter, req *http.Request, routers []*router.Router) {
	var (
		contenttype string
		content     string
		statuscode  int
	)
	switch r := res.(type) {
	case response.JSONResponse:
		statuscode = resolveStatusCode(r.StatusCode)
		contenttype = "application/json"
		contentBytes, err := json.Marshal(r.Content)
		if err != nil {
			content, statuscode = r.ResponseError(err.Error())
			break
		}
		content = string(contentBytes)
	case response.HTMLResponse:
		statuscode = resolveStatusCode(r.StatusCode)
		contenttype = "text/html"
		content = r.Content
	case response.Response:
		statuscode = resolveStatusCode(r.StatusCode)
		contenttype = "text/plain"
		content = r.Content
	}

	w.WriteHeader(statuscode)
	w.Header().Add("Content-Type", contenttype)
	fmt.Fprint(w, content)
}
