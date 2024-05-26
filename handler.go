package puff

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/nikumar1206/puff/middleware"
)

var fileContentTypes = map[string]string{
	"aac":    "audio/aac",
	"abw":    "application/x-abiword",
	"apng":   "image/apng",
	"arc":    "application/x-freearc",
	"avif":   "image/avif",
	"avi":    "video/x-msvideo",
	"azw":    "application/vnd.amazon.ebook",
	"bin":    "application/octet-stream",
	"bmp":    "image/bmp",
	"bz":     "application/x-bzip",
	"bz2":    "application/x-bzip2",
	"cda":    "application/x-cdf",
	"csh":    "application/x-csh",
	"css":    "text/css",
	"csv":    "text/csv",
	"doc":    "application/msword",
	"docx":   "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"eot":    "application/vnd.ms-fontobject",
	"epub":   "application/epub+zip",
	"gz":     "application/gzip",
	"gif":    "image/gif",
	"htm":    "text/html",
	"html":   "text/html",
	"ico":    "image/vnd.microsoft.icon",
	"ics":    "text/calendar",
	"jar":    "application/java-archive",
	"jpeg":   "image/jpeg",
	"jpg":    "image/jpeg",
	"js":     "text/javascript",
	"json":   "application/json",
	"jsonld": "application/ld+json",
	"mid":    "audio/midi",
	"midi":   "audio/midi",
	"mjs":    "text/javascript",
	"mp3":    "audio/mpeg",
	"mp4":    "video/mp4",
	"mpeg":   "video/mpeg",
	"mpkg":   "application/vnd.apple.installer+xml",
	"odp":    "application/vnd.oasis.opendocument.presentation",
	"ods":    "application/vnd.oasis.opendocument.spreadsheet",
	"odt":    "application/vnd.oasis.opendocument.text",
	"oga":    "audio/ogg",
	"ogv":    "video/ogg",
	"ogx":    "application/ogg",
	"opus":   "audio/opus",
	"otf":    "font/otf",
	"png":    "image/png",
	"pdf":    "application/pdf",
	"php":    "application/x-httpd-php",
	"ppt":    "application/vnd.ms-powerpoint",
	"pptx":   "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"rar":    "application/vnd.rar",
	"rtf":    "application/rtf",
	"sh":     "application/x-sh",
	"svg":    "image/svg+xml",
	"tar":    "application/x-tar",
	"tif":    "image/tiff",
	"tiff":   "image/tiff",
	"ts":     "video/mp2t",
	"ttf":    "font/ttf",
	"txt":    "text/plain",
	"vsd":    "application/vnd.visio",
	"wav":    "audio/wav",
	"weba":   "audio/webm",
	"webm":   "video/webm",
	"webp":   "image/webp",
	"woff":   "font/woff",
	"woff2":  "font/woff2",
	"xhtml":  "application/xhtml+xml",
	"xls":    "application/vnd.ms-excel",
	"xlsx":   "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"xml":    "application/xml",
	"xul":    "application/vnd.mozilla.xul+xml",
	"zip":    "application/zip",
	"3gp":    "video/3gpp",
	"3g2":    "video/3gpp2",
	"7z":     "application/x-7z-compressed",
}

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

func contentTypeFromFileSuffix(suffix string) string {
	ct := fileContentTypes[suffix]
	if ct == "" {
		return "text/plain" //we dont know the content type from file suffix
	}
	return ct
}

func Handler(w http.ResponseWriter, req *http.Request, route *Route) {
	defer func() {
		a := recover()
		if a != nil {
			errorID := middleware.RandomLogID()
			w.WriteHeader(500)
			w.Header().Add("Content-Type", "text/plain")
			fmt.Fprint(w, "There was a panic during the execution recovered by the handler. Error ID: "+errorID)
			slog.Error("Panic During Execution", slog.String("ERROR ID", errorID), slog.String("Error", a.(string)))
		}
	}()
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
		fileNameSplit := strings.Split(r.FileName, ".")
		suffix := fileNameSplit[len(fileNameSplit)-1]
		contentType = contentTypeFromFileSuffix(suffix)
		file, err := os.ReadFile(r.FileName)
		if err != nil {
			statusCode = 500
			content = "There was an error retrieving the file: " + err.Error()
		}
		statusCode = resolveStatusCode(r.StatusCode, req.Method)
		content = string(file)
	case Response:
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
