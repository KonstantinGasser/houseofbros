package user

import (
	"encoding/json"
	"log"

	"github.com/KonstantinGasser/houseofbros/socket"
	"github.com/gorilla/websocket"
)

type User struct {
	// Conn      *websocket.Conn
	Username  string        `json:"username"`
	Action    string        `json:"action"`
	Note      string        `json:"note"`
	Emojies   []interface{} `json:"emojies"`
	Reactions []interface{} `json:"reactions"`
	IsOnline  bool          `json:"is-online"`
}

func (u *User) Update(action, note string, emojies []interface{}) {
	u.Action = action
	u.Note = note
	u.Emojies = emojies
}

func (u *User) AddReaction(v interface{}) {
	u.Reactions = append(u.Reactions, v)
}

// runs in its own goroutine
func (user *User) ping(broadcast chan socket.Event, remove chan string, conn *websocket.Conn, uname string) {
	log.Printf("[started] ping to client <%s> goroutine\n", uname)
	defer func() {
		user.IsOnline = false
		conn.Close()
		b, _ := user.Serialize()
		broadcast <- socket.EventUser{
			Type: "user-changed",
			User: b,
		}
		remove <- uname
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			log.Printf("[ping] connection of <%s> failed to responed: %s\n", uname, err.Error())
			return
		}

	}
}

func (u *User) Serialize() ([]byte, error) {

	b, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return b, nil
}
