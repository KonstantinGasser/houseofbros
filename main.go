package main

import (
	"log"

	"github.com/KonstantinGasser/houseofbros/api"
)

func main() {
	server := api.NewServer()
	server.SetUp()

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
