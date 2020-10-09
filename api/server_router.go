package api

import "net/http"

type APIRouter struct{}

func (router *APIRouter) HandleFunc(path string, handlerFunc func(w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc(path, handlerFunc)
}
