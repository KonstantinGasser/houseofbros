package api

import (
	"log"
	"net/http"

	// init routes from package routes
	_ "github.com/KonstantinGasser/houseofbros/api/routes"
)

// HttpServer some docs
type HttpServer struct {
	addr string
}

// NewHTTPServer creates a new instance of a HTTPServer.
// Returns the pointer
func NewHTTPServer(addr string) *HttpServer {
	return &HttpServer{addr: addr}
}

// Serve spins up the server
func (server *HttpServer) Serve() error {
	log.Printf("[created] new Server started on %s\n", server.addr)
	if err := http.ListenAndServe(server.addr, nil); err != nil {
		return err
	}
	return nil
}
