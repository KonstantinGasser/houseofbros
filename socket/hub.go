package socket

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	tickerTime = 75 * time.Second
	// for upgrader
	readBuf  = 1024
	writeBuf = 1024
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  readBuf,
		WriteBufferSize: writeBuf,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

type Hub struct {
	Unknown chan *user
	Known   chan *user

	// protecting pool
	mu   sync.Mutex
	pool map[string]*user
}

// run runs in its own goroutine
func (hub *Hub) run() {
	ticker := time.NewTicker(tickerTime)
	defer func() {
		ticker.Stop()
		close(hub.Unknown)
		close(hub.Known)
	}()

	for {
		select {
		case user := <-hub.Unknown:
			log.Printf("[spooky] unknown user <%s> wants to connect ~ buh", user.username)
			if _, ok := hub.pool[user.username]; !ok {
				hub.mu.Lock()
				hub.pool[user.username] = user
				hub.mu.Unlock()
			}
			// tell the bros
			byteUser, _ := user.serialize() // TODO: Impl error chan??
			hub.stream(byteUser)

		case user := <-hub.Known:
			log.Printf("[cool-mate] nice to have you back <%s>", user.username)
			hub.mu.Lock()
			hub.pool[user.username] = user
			hub.mu.Unlock()

			// tell the bros
			byteUser, _ := user.serialize() // TODO: Impl error chan??
			hub.stream(byteUser)

		case <-ticker.C:
			log.Printf("[status] running: run() with %d goroutines, %d users", runtime.NumGoroutine(), len(hub.pool))
		default:
			// no requests on channels
		}
	}
}

func (hub *Hub) Upgrade(w http.ResponseWriter, r *http.Request, username string) error {
	log.Printf("[processing] upgrade request: %v", r.URL.Hostname())
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[shit] hub.Upgrade(): %v", err)
		return err
	}

	_user, ok := hub.pool[username]
	if ok {
		hub.mu.Lock()
		_user.conn = conn
		_user.isOnline = true
		hub.mu.Unlock()
		hub.Known <- _user

		// spin up user checkAlive
		go _user.checkAlive()
		return nil
	}

	user := &user{
		mu:       sync.Mutex{},
		hub:      hub,
		conn:     conn,
		username: username,
		isOnline: true,
		state:    createFuncyState(),
	}
	hub.Unknown <- user
	// as spin up user checkAlive
	go user.checkAlive()
	return nil
}

func (hub *Hub) UpdateUser(uname, action, note string, emo []interface{}) error {

	user, ok := hub.pool[uname]
	if !ok {
		return fmt.Errorf("dude there is no user for <%s>", uname)
	}
	_user := user.update(action, note, emo)
	byteUser, _ := _user.serialize()
	hub.stream(byteUser)

	return nil
}

func (hub *Hub) Reacte(reciever string, emojiType interface{}) error {
	user, ok := hub.pool[reciever]
	if !ok {
		log.Printf("[man-oh-man] user not found: %s", reciever)
		return fmt.Errorf("no user found: %s", reciever)
	}

	user.mu.Lock()
	defer user.mu.Unlock()
	user.state.reactions = append(user.state.reactions, emojiType)

	byteUser, _ := user.serialize()
	hub.stream(byteUser)
	log.Printf("[nice] reacted to: %s", reciever)
	return nil
}

func (hub *Hub) stream(msg []byte) {
	for _, u := range hub.pool {
		if !u.isOnline {
			continue
		}
		u.brocast(msg)
	}
}

func NewHub() *Hub {
	log.Println("[created] new SocketHub")
	hub := &Hub{
		Unknown: make(chan *user),
		Known:   make(chan *user),
		mu:      sync.Mutex{},
		pool:    make(map[string]*user),
	}
	go hub.run()
	return hub
}
