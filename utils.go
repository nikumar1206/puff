package puff

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"mime"
	"net/http"
	"strings"
)

func RandomNanoID() string {
	id := ""
	for range 4 {
		r := rand.IntN(25) + 1
		id += fmt.Sprintf("%c", ('A' - 1 + r))
	}
	id += "-"
	for range 4 {
		r := rand.IntN(9)
		id += fmt.Sprint(r)
	}
	return id
}

func resolveContentType(provided, default_content_type string) string {
	if provided == "" {
		return default_content_type
	}
	return provided
}

func resolveStatusCode(sc int, method string) int {
	if sc == 0 {
		if method == http.MethodPost {
			return http.StatusCreated
		}
		return http.StatusOK
	}
	return sc
}

func contentTypeFromFileName(name string) string {
	fileNameSplit := strings.Split(name, ".")
	suffix := fileNameSplit[len(fileNameSplit)-1]
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
