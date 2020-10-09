package api

import "net/http"

func (s *Server) HandlerProtcollUpgrade(w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		_json, err := decode(r.Body)
		if err != nil {
			http.Error(w, "`{'status': 404, 'message': 'body_to_json_decoding_error'}`", http.StatusBadRequest)
			return
		}
		if err := s.userHub.Create(w, r, _json); err != nil {
			http.Error(w, "`{'status': 500, 'message': 'create_user_upgrade_protocoll'}`", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("`{'status': 201, 'message': 'user_created_protocoll_upgraded'}`"))
	}
}
