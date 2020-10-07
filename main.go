package main

import (
	"log"

	"github.com/KonstantinGasser/houseofbros/api"
)

func main() {
	server := api.NewServer()
	server.Routes()
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
