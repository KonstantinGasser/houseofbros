package websocket

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/KonstantinGasser/houseofbros/status"
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
	Join    chan *connection
	Remove  chan string
	Brocast chan map[string]interface{}
	// mutex for following
	mu   sync.Mutex
	bros map[string]*connection
}

type connection struct {
	hub *SocketHub
	mu  sync.Mutex
	c   *websocket.Conn
	bro *status.Bro
}

// run runs in its own goroutine
func (hub *SocketHub) run() {
	ticker := time.NewTicker(tickerTime)
	defer ticker.Stop()

	for {
		select {
		case c := <-hub.Join: //add new connection to hub
			hub.mu.Lock()
			if _, ok := hub.bros[c.bro.Uname]; !ok {
				hub.bros[c.bro.Uname] = c
			}
			hub.mu.Unlock()
			log.Printf("[appended] connection <%s> added to connections\n", c.bro.Uname)

			for _, bro := range hub.bros { // update: new bro incoming
				data := map[string]interface{}{
					"uname":    c.bro.Uname,    //bro.bro.Uname,
					"activity": c.bro.Activity, //bro.bro.Activity,
					"comment":  c.bro.Comment,  //bro.bro.Comment,
					"emotion":  c.bro.Emotion,  //bro.bro.Emotion,
				}
				b, err := json.Marshal(data)
				if err != nil {
					log.Printf("[error] json.Marshal(data): %v", err.Error())
					continue
				}
				if err := connSend(bro, b); err != nil {
					continue
				}
			}
		case uname := <-hub.Remove: // delete connection from hub
			hub.mu.Lock()
			delete(hub.bros, uname)
			hub.mu.Unlock()
			log.Printf("[removed] user <%s> removed from connections", uname)
		case bro := <-hub.Brocast: // broadcast status update
			for _, br := range hub.bros {
				fmt.Println(br.bro.Uname)
				_bro := hub.bros[bro["uname"].(string)].bro.UpdateBro(
					bro["activity"].(string),
					bro["comment"].(string),
					bro["emotion"].([]interface{}),
				)
				b, err := json.Marshal(map[string]interface{}{
					"uname":    _bro.Uname,
					"activity": _bro.Activity,
					"comment":  _bro.Comment,
					"emotion":  _bro.Emotion,
				})
				if err != nil {
					continue
				}
				if err := connSend(br, b); err != nil {
					hub.Remove <- br.bro.Uname
					continue
				}
			}
		case <-ticker.C:
			log.Printf("[status] SocketHub - running: current connections: %d\n", len(hub.bros))
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
	if _, ok := hub.bros[uname]; ok {
		log.Printf("[error] uname: %s already exists", uname)
		return fmt.Errorf("`{'error': 'username-existis-exception'}`")
	}

	connection := &connection{
		hub: hub,
		mu:  sync.Mutex{},
		c:   conn,
		bro: status.NewBro(uname),
	}
	hub.bros[uname] = connection

	hub.Join <- connection
	return nil
}

func parseURL(r *http.Request) (uname string) {
	uname = r.URL.Query().Get("uname")
	return uname
}

// BroYouThere checks if there is a connection for a bro
func (hub *SocketHub) BroYouThere(uname string) bool {
	if _, ok := hub.bros[uname]; ok {
		return true
	}
	return false
}

// NewHub creates a new SocketHub and starts the run method
// returns a pointer to the SocketHub
func NewHub() *SocketHub {
	log.Print("[created] new SocketHub\n")
	socket := SocketHub{
		Join:    make(chan *connection),
		Remove:  make(chan string),
		Brocast: make(chan map[string]interface{}),
		mu:      sync.Mutex{},
		bros:    make(map[string]*connection),
	}

	// spin-up goroutine
	go socket.run()
	return &socket
}

func connSend(conn *connection, msg []byte) error {
	log.Printf("SEND: %s\n", string(msg))
	var w io.WriteCloser
	conn.mu.Lock()
	defer conn.mu.Unlock()
	w, err := conn.c.NextWriter(websocket.TextMessage)
	if err != nil {
		// clean-up
		log.Printf("[error] conn.c.NextWriter(...) unable to create NextWriter: %v", err.Error())
		conn.hub.Remove <- conn.bro.Uname
		connClose(conn)
		return err
	}
	defer w.Close()
	if _, err := w.Write(msg); err != nil {
		// clean-up
		log.Printf("[error] w.Write(msg) unable to send message to connection <%v> :%v", &conn.c, err.Error())
		conn.hub.Remove <- conn.bro.Uname
		connClose(conn)
		return err
	}
	log.Printf("[send] message send to <%v>", &conn.c)
	return nil
}

// listenForClose runs as a goroutine for each connection
// trying to listen for the client - if that fails
func listenForClose(conn *connection) {

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

func (hub *SocketHub) DecodeFullMap() ([]byte, error) {
	var payload = make(map[string]interface{})
	for _, bro := range hub.bros {
		payload[bro.bro.Uname] = map[string]interface{}{
			"uname":    bro.bro.Uname,
			"activity": bro.bro.Activity,
			"comment":  bro.bro.Comment,
			"emotion":  bro.bro.Emotion,
		}
	}

	_data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return _data, nil
}
