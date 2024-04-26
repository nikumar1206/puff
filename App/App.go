package app

import (
	"fmt"
	"net/http"
	handler "puff/Handler"
	route "puff/Route"
	router "puff/Router"
)

type AppI struct {
	Network bool //host to the entire network?
	Reload  bool //live reload?
	Port    int  //port number to use
	Routes  []route.Route
	Routers []*router.Router
	// Middlewares
}

func (ac *AppI) IncludeRouter(r *router.Router) {
	ac.Routers = append(ac.Routers, r)
}

func (ac *AppI) sendToHandler(w http.ResponseWriter, req *http.Request) {
	handler.Handler(w, req, ac.Routers)
}

func (ac *AppI) ListenAndServe() {
	http.HandleFunc("", ac.sendToHandler)
	network := ""
	if ac.Network {
		network += "0.0.0.0"
	}
	network += ":" + fmt.Sprintf("%d", ac.Port)
	http.ListenAndServe(network, nil)
}
