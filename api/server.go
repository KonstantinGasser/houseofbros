package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/KonstantinGasser/houseofbros/socket"
)

type Server struct {
	hub    *socket.Hub
	router *router
}

func NewServer() *Server {
	log.Print("[created] new Server started on :8080\n")
	return &Server{
		hub:    socket.NewHub(),
		router: &router{},
	}
}

func (s *Server) Routes() {
	s.router.HandleFunc("/api/v1/websocket/connect", s.HandleSocketConnection())
	s.router.HandleFunc("/api/v1/bros/update", s.HandelUpdate())
	s.router.HandleFunc("/api/v1/bros/reaction", s.HandelReaction())
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(":8080", nil)
}

func decode(body io.ReadCloser) (map[string]interface{}, error) {
	defer body.Close()

	var data map[string]interface{}

	if err := json.NewDecoder(body).Decode(&data); err != nil {
		log.Printf("[oh-man] decode(): %v", err.Error())
		return nil, err
	}
	return data, nil
}
