package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	response "puff/Response"
	router "puff/Router"
)

// func UserHandler() interface{} { //this is an example handler that our handler function is looking for
// 	return response.JSONResponse{
// 		StatusCode: 200, //of course this is optional (default: 200)
// 		Content:    map[string]interface{}{"food": "pizza", "taste": "yummy"},
// 	}
// }

func statusCode(sc int) int {
	if sc == 0 {
		return 200
	} else {
		return sc
	}
}

func Handler(w http.ResponseWriter, req *http.Request, routers []*router.Router) {
	var contenttype string
	var content string
	var statuscode int
	switch r := res.(type) {
	case response.JSONResponse:
		statuscode = statusCode(r.StatusCode)
		contenttype = "application/json"
		contentBytes, err := json.Marshal(r.Content)
		if err != nil {
			content, statuscode = r.ResponseError(err.Error())
			break
		}
		content = string(contentBytes)
	case response.HTMLResponse:
		statuscode = statusCode(r.StatusCode)
		contenttype = "text/html"
		content = r.Content
	case response.Response:
		statuscode = statusCode(r.StatusCode)
		contenttype = "text/plain"
		content = r.Content
	}

	w.WriteHeader(statuscode)
	w.Header().Add("Content-Type", contenttype)
	fmt.Fprint(w, content)
}
