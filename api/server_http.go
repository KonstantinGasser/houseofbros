package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
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
	mainHub := socket.NewMainHub()
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
	// routes
}

type APIRouter struct{}

func (router *APIRouter) HandleFunc(path string, handlerFunc func(w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc(path, handlerFunc)
}

func decode(body io.ReadCloser) (map[string]interface{}, error) {

	var data map[string]interface{}
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return data, nil
}
