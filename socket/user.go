package socket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pingPeriod = 3 * time.Second
)

type user struct {
	username string
	isOnline bool
	state    *state
	// protects conn, hub
	mu   sync.Mutex
	conn *websocket.Conn
	hub  *Hub
}

type state struct {
	action    string
	note      string
	emojies   []interface{}
	reactions []interface{}
}

func (user *user) brocast(msg []byte) {
	wr, err := user.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		log.Printf("[dude!] user.brocast(msg).NextWriter: %v\n", err)
		return
	}
	defer wr.Close()

	if _, err := wr.Write(msg); err != nil {
		log.Printf("[dude!!] user.brocast(msg).Write: %v\n", err)
		return
	}
}

// checkAlive runs in its own goroutine for each user
func (user *user) checkAlive() {
	defer func() {
		user.isOnline = false
		user.conn.Close()
		b, _ := user.serialize()
		user.hub.stream(b)
	}()
	for {
		// time.Sleep(pingPeriod)

		_, _, err := user.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[dude-you-there?] connection broke to: %s", user.username)
			}
			break
		}
	}
}

func (user *user) update(action, note string, emo []interface{}) *user {
	user.mu.Lock()
	defer user.mu.Unlock()
	user.state.action = action
	user.state.note = note
	user.state.emojies = emo
	user.state.reactions = []interface{}{}
	return user
}

func (user *user) serialize() ([]byte, error) {

	var u = map[string]interface{}{
		"username":  user.username,
		"is-online": user.isOnline,
		"action":    user.state.action,
		"note":      user.state.note,
		"emojies":   user.state.emojies,
		"reactions": user.state.reactions,
	}
	b, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	log.Printf("[watch-out-bro] to: %s\n", user.username)
	return b, nil
}

func createFuncyState() *state {
	return &state{
		action:    "some action",
		note:      "how much is the fish?",
		emojies:   []interface{}{2, 4},
		reactions: []interface{}{},
	}
}
