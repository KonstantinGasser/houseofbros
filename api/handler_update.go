package api

import (
	"log"
	"net/http"
)

func (s *Server) HandelUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[called] route: %v", r.URL.Path)
		_json, err := decode(r.Body)
		if err != nil {
			http.Error(w, "sorry mate ~ your json is just not quite as we like it", http.StatusBadRequest)
			return
		}
		if err := s.hub.UpdateUser(
			_json["username"].(string),
			_json["action"].(string),
			_json["note"].(string),
			_json["emojies"].([]interface{})); err != nil {
			http.Error(w, "sorry mate we messed up this request", http.StatusInternalServerError)
		}
	}
}

func (s *Server) HandelReaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[called] route: %v", r.URL.Path)
		_json, err := decode(r.Body)
		if err != nil {
			http.Error(w, "sorry mate ~ your json is just not quite as we like it", http.StatusBadRequest)
			return
		}
		if err := s.hub.Reacte(_json["for"].(string), _json["with"]); err != nil {
			http.Error(w, "sorry mate ~ your json is just not quite as we like it", http.StatusBadRequest)
			return
		}
	}
}
