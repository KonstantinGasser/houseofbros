package user

import (
	"log"
	"net/http"
	"sync"

	"github.com/KonstantinGasser/houseofbros/socket"
)

type UserHub struct {
	mainHub *socket.MainHub
	mu      sync.Mutex
	users   map[string]User `json:"users"`
}

func (hub *UserHub) Create(w http.ResponseWriter, r *http.Request, v map[string]interface{}) error {

	uname := v["uname"].(string)
	user := &StdUser{
		Username:  uname,
		Action:    "Beaking the Desk",
		Note:      "How much is the fish?",
		Emojies:   []interface{}{1},
		Reactions: []interface{}{},
	}
	hub.mu.Lock()
	hub.users[uname] = user
	hub.mu.Unlock()

	log.Printf("[created] new User <%s>\n", uname)
	if err := hub.mainHub.Upgrade(w, r, uname); err != nil {
		return err
	}
	return nil

}

func (hub *UserHub) Update(v map[string]interface{}) error {
	return nil
}

func (hub *UserHub) Delete(v map[string]interface{}) error {
	return nil
}

func (hub *UserHub) Serialize() ([]byte, error) {
	return nil, nil
}
