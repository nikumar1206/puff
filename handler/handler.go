package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"puff/field"
	"puff/request"
	response "puff/response"
	"puff/route"
	"strconv"
)

func resolveStatusCode(sc int) int {
	if sc == 0 {
		return 200
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
	var fields map[string]interface{}
	fields = make(map[string]interface{})
	for _, routeField := range route.Fields {
		var value string
		if req.Method == "GET" {
			value = req.PathValue(routeField.Name)
		} else if req.Method == "POST" || req.Method == "PUT" || req.Method == "PATH" { //these are the only methods that allow for req.PostForm
			value = req.PostForm.Get(routeField.Name)
		}
		if value == "" {
			http.Error(w, routeField.MissingFieldError(), 422)
		}
		if !routeField.Validate(value) {
			http.Error(w, routeField.TypeValidationError(), 422)
			return
		}
		var typedVal interface{} = value
		pstt := field.ParseStringToType(value)
		switch pstt {
		case "int":
			typedVal, _ = strconv.Atoi(value)
		case "bool":
			typedVal, _ = strconv.ParseBool(value)
		}
		fields[routeField.Name] = typedVal
	}
	requestDetails.Fields = fields
	res := route.Handler(requestDetails) // FIX ME: we should give the user handle function a request body as well
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
			content = r.ResponseError(err)
			http.Error(w, content, 500)
		}
		content = string(contentBytes)
	case response.HTMLResponse:
		statusCode = resolveStatusCode(r.StatusCode)
		contentType = "text/html"
		content = r.Content
	case response.Response:
		statusCode = resolveStatusCode(r.StatusCode)
		contentType = resolveContentType(r.ContentType)
		content = r.Content
	}

	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", contentType)
	fmt.Fprint(w, content)
}
