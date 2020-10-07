package api

import (
	"net/http"
)

type router struct{}

func (router *router) HandleFunc(path string, fn func(w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc(path, fn)
}
