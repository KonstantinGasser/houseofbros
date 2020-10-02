package websocket

import (
	"fmt"
	"log"
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

// SocketHub knows all the connected users and updates them
// if either changes its status
type SocketHub struct {
	Join      chan map[string]interface{}
	Remove    chan string
	Broadcast chan map[string]interface{}
	// mutex for following
	mu    sync.Mutex
	conns map[string]*connection
}

type connection struct {
	hub   *SocketHub
	uname string
	mu    sync.Mutex
	c     *websocket.Conn
}

// run runs in its own goroutine
func (hub *SocketHub) run() {
	ticker := time.NewTicker(tickerTime)
	defer ticker.Stop()

	for {
		select {
		case c := <-hub.Join: //add new connection to hub
			hub.mu.Lock()
			if _, ok := hub.conns[c["uname"].(string)]; !ok {
				hub.conns[c["uname"].(string)] = c["conn"].(*connection)
			}
			hub.mu.Unlock()
			log.Printf("[appended] connection <%s:%v> added to connections", c["uname"], c["conn"])
		case uname := <-hub.Remove: // delete connection from hub
			hub.mu.Lock()
			delete(hub.conns, uname)
			hub.mu.Unlock()
			log.Printf("[removed] user <%s> removed from connections", uname)
		case update := <-hub.Broadcast: // broadcast status update
			log.Println(update)

			connSend(update["conn"].(*connection), update["msg"].([]byte))
		case <-ticker.C:
			log.Print("[status] SocketHub - running\n")
		}
	}

}

// UpgradeServe takes the incoming HTTP requests and upgrades the protocoll
// to be a WS one
func (hub *SocketHub) UpgradeServe(w http.ResponseWriter, r *http.Request) error {

	uname := parseURL(r)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[error] upgrader.Upgrade for: %v\n", r)
		return err
	}

	hub.mu.Lock()
	defer hub.mu.Unlock()
	if _, ok := hub.conns[uname]; ok {
		log.Printf("[error] uname: %s already exists", uname)
		return fmt.Errorf("`{'error': 'username-existis-exception'}`")
	}
	log.Printf("[appended] connection <%s:%v> added to connections\n", uname, &conn)
	connection := &connection{
		hub:   hub,
		uname: uname,
		mu:    sync.Mutex{},
		c:     conn,
	}
	hub.conns[uname] = connection
	// let connection listen
	go connRead(connection)

	// send ACK to connection
	if err := connSend(connection, []byte("Hello Bro")); err != nil {
		log.Printf("[error] sending ack to connection: %v", &connection.c)
		return err
	}
	return nil
}

func parseURL(r *http.Request) (uname string) {
	uname = r.URL.Query().Get("uname")
	return uname
}

// NewHub creates a new SocketHub and starts the run method
// returns a pointer to the SocketHub
func NewHub() *SocketHub {
	log.Print("[created] new SocketHub\n")
	socket := SocketHub{
		Join:      make(chan map[string]interface{}),
		Remove:    make(chan string),
		Broadcast: make(chan map[string]interface{}),
		mu:        sync.Mutex{},
		conns:     make(map[string]*connection),
	}

	// spin-up goroutine
	go socket.run()
	return &socket
}

// connRead runs in its own goroutine
func connRead(conn *connection) {
	defer func() {
		conn.hub.Remove <- conn.uname
	}()
	for {
		_, msg, err := conn.c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[error] connRead/UnexcpectedCloseError: %v", err.Error())
			}
			log.Printf("[error] conn.c.ReadMessage(): %v", err.Error())
			break
		}
		conn.hub.Broadcast <- map[string]interface{}{
			"conn": conn,
			"msg":  msg,
		}
	}
}

// writes
func connSend(conn *connection, msg []byte) error {

	conn.mu.Lock()
	defer conn.mu.Unlock()
	w, err := conn.c.NextWriter(websocket.TextMessage)
	if err != nil {
		log.Printf("[error] conn.c.NextWriter(...) unable to create NextWriter: %v", err.Error())
		return err
	}
	if _, err := w.Write(msg); err != nil {
		log.Printf("[error] w.Write(msg) unable to send message to connection <%v> :%v", &conn.c, err.Error())
		return err
	}
	log.Printf("[send] message send to <%v>", &conn.c)
	return nil
}

// connClose send closure message to client
func connClose(conn *connection) error {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	if err := conn.c.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
		return err
	}
	return nil
}
