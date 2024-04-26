package handler

import (
	"net/http"
	router "puff/Router"
)

func Handler(w http.ResponseWriter, req *http.Request, routers []*router.Router) { //should probably add middleware paramaters here
	//you want to handle middlewares inside this functions
}
