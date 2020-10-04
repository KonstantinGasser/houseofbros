package routes

import (
	"log"
	"net/http"
)

// forwardUpgrade sends the request to the websocket for registration
func routeUpgrade(w http.ResponseWriter, r *http.Request) {
	if err := hub.UpgradeServe(w, r); err != nil {
		log.Printf("[error] protocoll failed to upgrade: %v", err.Error())
		// fmt.Fprintf(w, "Protocoll failed to upgrade: %v", err.Error())
	}
}
