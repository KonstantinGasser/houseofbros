package main

import (
	"log"

	"github.com/KonstantinGasser/houseofbros/api"
)

func main() {

	server := api.NewHTTPServer(":8080")

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
