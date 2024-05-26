package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nikumar1206/puff/utils"
)

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			a := recover()
			if a != nil {
				errorID := utils.RandomNanoID()
				w.WriteHeader(500)
				w.Header().Add("Content-Type", "text/plain")
				fmt.Fprint(w, "There was a panic during the execution recovered by the panic handling middleware. Error ID: "+errorID)
				slog.Error("Panic During Execution", slog.String("ERROR ID", errorID), slog.String("Error", a.(string)))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
