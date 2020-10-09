package socket

import (
	"net/http"
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

type MainHub struct {
	Singlecast chan []byte
	Broadcast  chan []byte
	// protects connection map
	mu    sync.Mutex
	conns map[string]*websocket.Conn
}

// run runs in its own goroutine
func (hub *MainHub) run() {
	ticker := time.NewTicker(tickerTime)
	defer ticker.Stop()

	for {
		select {
		case msg := <-hub.Broadcast:
			// do something here

		case msg := <-hub.Singlecast:
			// do something here

		case <-ticker.C:
			// log something
		default:
			// do nothing
		}
	}
}

func (hub *MainHub) Upgrade(w http.ResponseWriter, r *http.Request, username string) error {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	hub.mu.Lock()
	hub.conns[username] = conn
	hub.mu.Unlock()

	return nil
}

func NewMainHub() *MainHub {
	return &MainHub{
		Singlecast: make(chan []byte),
		Broadcast:  make(chan []byte),
		mu:         sync.Mutex{},
		conns:      make(map[string]*websocket.Conn),
	}
}
