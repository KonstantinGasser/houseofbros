package api

import "net/http"

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
	if err := http.ListenAndServe(server.addr, nil); err != nil {
		return err
	}
	return nil
}
