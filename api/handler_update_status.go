package api

import (
	"net/http"

	"github.com/KonstantinGasser/houseofbros/socket"
)

func (s *Server) HandlerUpdateStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Conent-Type", "application/json")
		_json, err := decode(r.Body)
		if err != nil {
			http.Error(w, "`{'status': 404, 'message': 'malformed_request_body'}`", http.StatusBadRequest)
			return
		}

		if err := s.userHub.UpdateStatus(_json); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		b, err := s.userHub.Serialize()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.mainHub.Broadcast <- socket.EventUser{
			Type: "user-new",
			User: b,
		}
		w.WriteHeader(http.StatusOK)
	}
}
