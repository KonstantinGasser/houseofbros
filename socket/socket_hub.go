package socket

import (
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	tickerTime = 5 * time.Second
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
	Singlecast chan Event //[]byte
	Broadcast  chan Event //[]byte
	Remove     chan string
	// protects connection map
	mu    sync.Mutex
	conns map[string]*websocket.Conn
}

// Run runs in its own goroutine
func (hub *MainHub) Run() {
	ticker := time.NewTicker(tickerTime)
	defer ticker.Stop()

	for {
		select {
		case msg := <-hub.Broadcast:
			// do something here
			b, _ := msg.Serialize()
			stream(b, hub.conns)
		case <-hub.Singlecast:
			// do something here
		case uname := <-hub.Remove:
			log.Printf("[remove] user <%s>", uname)
			hub.mu.Lock()
			delete(hub.conns, uname)
			hub.mu.Unlock()
		case <-ticker.C:
			log.Printf("[status] run(): running, goroutines: %d, connections: %d", runtime.NumGoroutine(), len(hub.conns))
			// log something
		default:
			// do nothing
		}
	}
}

func stream(msg []byte, cs map[string]*websocket.Conn) error {
	var _err error
	for _, c := range cs {
		wr, err := c.NextWriter(websocket.TextMessage)
		if err != nil {
			_err = err
			continue
		}
		wr.Write(msg)
		wr.Close()
	}
	return _err
}

func (hub *MainHub) Upgrade(w http.ResponseWriter, r *http.Request, username string) (*websocket.Conn, error) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	hub.mu.Lock()
	hub.conns[username] = conn
	hub.mu.Unlock()

	return conn, nil
}

func NewMainHub() *MainHub {
	log.Printf("[created] new MainHub\n")
	return &MainHub{
		Singlecast: make(chan Event),
		Broadcast:  make(chan Event),
		Remove:     make(chan string),
		mu:         sync.Mutex{},
		conns:      make(map[string]*websocket.Conn),
	}
}
