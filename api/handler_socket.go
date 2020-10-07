package api

import (
	"log"
	"net/http"
)

// HandleSocketConnection some docs
func (s *Server) HandleSocketConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[called] route: %v", r.URL.Path)
		uname := r.URL.Query().Get("uname")
		if len(uname) == 0 {
			http.Error(w, "sorry mate that name is just not what we are looking for", http.StatusBadRequest)
			return
		}

		if err := s.hub.Upgrade(w, r, uname); err != nil {
			http.Error(w, "sorry mate", http.StatusInternalServerError)
			return
		}
	}
}
