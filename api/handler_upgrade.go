package api

import (
	"net/http"

	"github.com/KonstantinGasser/houseofbros/socket"
)

func (s *Server) HandlerProtcollUpgrade() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		uname := r.URL.Query().Get("uname")
		b, err := s.userHub.Create(w, r, uname)
		if err != nil {
			http.Error(w, "`{'status': 500, 'message': 'create_user_upgrade_protocoll'}`", http.StatusInternalServerError)
			return
		}
		// b, _ := s.userHub.Serialize()
		s.mainHub.Broadcast <- socket.EventUser{
			Type: "user-new",
			User: b,
		}
	}
}
