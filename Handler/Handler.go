package handler

import (
	"fmt"
	"net/http"
	router "puff/Router"
)

func Handler(
	w http.ResponseWriter, req *http.Request, routers []*router.Router,
) {
	fmt.Fprint(w, "Hello from Puff 💨")
}
