package api

import "net/http"

func (s *Server) HandlerProtcollUpgrade() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := s.userHub.Create(w, r); err != nil {
			http.Error(w, "`{'status': 500, 'message': 'create_user_upgrade_protocoll'}`", http.StatusInternalServerError)
			return
		}
	}
}
