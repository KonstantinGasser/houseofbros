package routes

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func routeGetCurrent(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "`{'status': 405, 'message': 'method not allowd'}`", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "Bro-Server since 99")

	_data, err := hub.DecodeFullMap()
	if err != nil {
		http.Error(w, "`{'status': 500, 'message': 'sorry house_of_bros is experiancing some issues..'}`", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(_data)

}

func routeUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "`{'status': 405,'message': 'method not allowd'}`", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "Bro-Server since 99")
	data, err := decode(r.Body)
	if err != nil {
		http.Error(w, "`{'status': 500, 'message': 'bronation is having troubles processing your request ~ sry'}`", http.StatusInternalServerError)
		return
	}
	if ok := hub.BroYouThere(data["uname"].(string)); !ok {
		http.Error(w, "{'status': 404, 'message': 'sorry mate - your are not part of the bronation yet'}", http.StatusBadRequest)
		return
	}
	hub.Brocast <- data
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("`{'status': 204, 'message': 'Your Bro-Mantra was succefully updated'}`"))

}

func decode(body io.ReadCloser) (map[string]interface{}, error) {
	defer body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		log.Printf("[error] unable to decode Request.Body: %v", err)
		return nil, err
	}
	return data, nil
}
