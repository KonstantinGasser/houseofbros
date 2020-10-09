package user

import (
	"net/http"
	"sync"

	"github.com/KonstantinGasser/houseofbros/socket"
)

type UserStorage interface {
	Create(w http.ResponseWriter, r *http.Request, v map[string]interface{}) error
	Update(v map[string]interface{}) error
	Delete(v map[string]interface{}) error
	UUID() (string, error)
	Serialize() ([]byte, error)
}

func NewUserHub(mainHub *socket.MainHub) UserStorage {
	return &UserHub{
		mainHub: mainHub,
		mu:      sync.Mutex{},
		users:   make(map[string]User),
	}
}
