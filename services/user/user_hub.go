package user

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/KonstantinGasser/houseofbros/socket"
)

// type User interface {
// 	Serialize() ([]byte, error)
// }

type UserHub struct {
	mainHub *socket.MainHub
	mu      sync.Mutex
	Users   map[string]*User `json:"users"`
}

func (hub *UserHub) Create(w http.ResponseWriter, r *http.Request, uname string) ([]byte, error) {

	conn, err := hub.mainHub.Upgrade(w, r, uname)
	if err != nil {
		return nil, err
	}

	_reqUser, ok := hub.Users[uname]
	if !ok {
		_reqUser = &User{
			Username:  uname,
			Action:    "Beaking the desk",
			Note:      "How much is the fish?",
			IsOnline:  true,
			Emojies:   []interface{}{0},
			Reactions: []interface{}{},
		}
		hub.mu.Lock()
		hub.Users[uname] = _reqUser
		hub.mu.Unlock()
	}
	_reqUser.IsOnline = true
	go _reqUser.ping(hub.mainHub.Broadcast, hub.mainHub.Remove, conn, uname)

	log.Printf("[created] User <%s>\n", uname)
	b, _ := _reqUser.Serialize()
	return b, nil

}

func (hub *UserHub) UpdateStatus(v map[string]interface{}) error {
	uname := v["username"].(string)
	action := v["action"].(string)
	note := v["note"].(string)
	emojies := v["emojies"].([]interface{})

	if _, ok := hub.Users[uname]; !ok {
		return fmt.Errorf("`{'status': 404, 'messasge': 'user_not_found'}`")
	}
	hub.mu.Lock()
	user, _ := hub.Users[uname]
	user.Update(action, note, emojies)
	hub.mu.Unlock()

	return nil
}

func (hub *UserHub) AddReaction(v map[string]interface{}) error {
	user, ok := hub.Users[v["to-user"].(string)]
	if !ok {
		return fmt.Errorf("`{'status': 404, 'messasge': 'user_not_found'}`")
	}
	user.AddReaction(v["reaction"])
	log.Printf("[added] reaction from:%v - to:%v\n", v["from-user"], v["to-user"])
	return nil
}

func (hub *UserHub) Delete(v map[string]interface{}) error {
	return nil
}

func (hub *UserHub) UUID() (string, error) {

	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	b[7] = 0x40
	b[9] = 0x80

	_hex := hex.EncodeToString(b)
	hexString := strings.ReplaceAll(_hex, " ", "")

	uuid := []string{hexString[:8], "-", hexString[8:12], "-", hexString[12:16], "=", hexString[16:20], "=", hexString[20:]}
	return strings.Join(uuid, ""), nil
}

func (hub *UserHub) Serialize() ([]byte, error) {

	b, err := json.Marshal(hub)
	if err != nil {
		log.Printf("[error] userHub.Serialize(): %v", err.Error())
		return nil, err
	}
	return b, nil
}
