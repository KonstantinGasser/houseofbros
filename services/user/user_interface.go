package user

import (
	"log"
	"net/http"
	"sync"

	"github.com/KonstantinGasser/houseofbros/socket"
)

type UserStorage interface {
	Create(w http.ResponseWriter, r *http.Request, uname string) ([]byte, error)
	UpdateStatus(v map[string]interface{}) error
	AddReaction(v map[string]interface{}) error
	Delete(v map[string]interface{}) error
	UUID() (string, error)
	Serialize() ([]byte, error)
}

func NewUserHub(mainHub *socket.MainHub) UserStorage {
	log.Printf("[created] new UserStorage\n")
	return &UserHub{
		mainHub: mainHub,
		mu:      sync.Mutex{},
		Users:   make(map[string]*User),
	}
}
