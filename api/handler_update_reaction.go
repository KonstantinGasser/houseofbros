package api

import (
	"net/http"

	"github.com/KonstantinGasser/houseofbros/socket"
)

func (s *Server) HandlerUpdateReaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_json, err := decode(r.Body)
		if err != nil {
			http.Error(w, "`{'status': 404, 'message': 'malformed_request_body'}`", http.StatusBadRequest)
			return
		}

		if err := s.userHub.AddReaction(_json); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.mainHub.Broadcast <- socket.EventReaction{
			Type:     "user-reaction",
			From:     _json["from-user"].(string),
			To:       _json["to-user"].(string),
			Reaction: _json["reaction"],
		}
		w.WriteHeader(http.StatusOK)

	}
}
