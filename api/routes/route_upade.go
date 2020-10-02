package routes

import (
	"fmt"
	"net/http"
)

func routeUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Route /status/update works!")
}
