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
	http.HandleFunc("/api/v1/status", routeStatus)
	log.Println("[set-up] route: /status")

	// upgrade and start websocket connection
	http.HandleFunc("/api/v1/websocket", routeUpgrade)
	log.Println("[set-up] route: /websocket")

	// route: update status
	http.HandleFunc("/api/v1/status/update", routeUpdate)
	log.Println("[set-up] route: /status/update")
}
