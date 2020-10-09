package user

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/KonstantinGasser/houseofbros/socket"
)

type User interface {
	Serialize() ([]byte, error)
}

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
		return nil, err
	}
	return b, nil
}
