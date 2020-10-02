package routes

import (
	"fmt"
	"net/http"
)

func routeStatus(w http.ResponseWriter, r *http.Request) {
	// get status of all connected people
	fmt.Fprintf(w, "Route /status works!")
}
