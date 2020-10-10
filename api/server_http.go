package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/KonstantinGasser/houseofbros/services/card"
	"github.com/KonstantinGasser/houseofbros/services/user"
	"github.com/KonstantinGasser/houseofbros/socket"
)

type Server struct {
	mainHub *socket.MainHub
	userHub user.UserStorage
	cardHub card.CardStorage

	router *APIRouter
}

func NewServer() *Server {
	log.Printf("[created] new API-Server\n")
	mainHub := socket.NewMainHub()
	go mainHub.Run()
	return &Server{
		mainHub: mainHub,
		userHub: user.NewUserHub(mainHub),
		cardHub: card.NewCardHub(mainHub),
	}
}

func (s *Server) Run() error {
	return http.ListenAndServe(":8080", nil)
}

func (s *Server) SetUp() {
	// routes: user
	log.Printf("[set-up] route: /api/v1/ws\n")
	s.router.HandleFunc("/api/v1/ws", s.HandlerProtcollUpgrade())

	log.Printf("[set-up] route: /api/v1/user/all\n")
	s.router.HandleFunc("/api/v1/user/all", s.HandlerUserAll())

	log.Printf("[set-up] route: /api/v1/user/update\n")
	s.router.HandleFunc("/api/v1/user/update", s.HandlerUpdateStatus())

	log.Printf("[set-up] route: /api/v1/user/reaction\n")
	s.router.HandleFunc("/api/v1/user/reaction", s.HandlerUpdateReaction())

	// route: card
	log.Printf("[set-up] route: /api/v1/card/all\n")
	s.router.HandleFunc("/api/v1/card/all", s.HandlerCardAll())
}

func decode(body io.ReadCloser) (map[string]interface{}, error) {

	var data map[string]interface{}
	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("[error] decode(): %s\n", err.Error())
		return nil, err
	}
	if err := json.Unmarshal(b, &data); err != nil {
		log.Printf("[error] decode(): %s\n", err.Error())
		return nil, err
	}
	return data, nil
}
