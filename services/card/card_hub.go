package card

import (
	"sync"

	"github.com/KonstantinGasser/houseofbros/socket"
)

type Card interface {
	Serialize() ([]byte, error)
}

type CardHub struct {
	mainHub *socket.MainHub
	mu      sync.Mutex
	Cards   map[string]*Card `json:"cards"`
}

func (hub *CardHub) Create() {}

func (hub *CardHub) Update() {}

func (hub *CardHub) Delete() {}

func (hub *CardHub) Serialize() ([]byte, error) {
	return nil, nil
}
