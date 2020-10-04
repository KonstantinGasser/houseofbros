package routes

import (
	"log"
	"net/http"

	"github.com/KonstantinGasser/houseofbros/websocket"
)

var hub *websocket.SocketHub

func init() {
	hub = websocket.NewHub()
	// route: status page
	http.HandleFunc("/api/v1/bros", routeStatus)
	log.Println("[set-up] route: /api/v1/bros/")

	// upgrade and start websocket connection
	http.HandleFunc("/api/v1/websocket", routeUpgrade)
	log.Println("[set-up] route: /api/v1/websocket")

	// route: update status
	http.HandleFunc("/api/v1/bros/update", routeUpdate)
	log.Println("[set-up] route: /api/v1/bros/update")

	// route: get current state
	http.HandleFunc("/api/v1/bros/current", routeGetCurrent)
	log.Println("[set-up] route: /api/v1/bros/current")
}
