package card

import (
	"sync"

	"github.com/KonstantinGasser/houseofbros/socket"
)

type CardStorage interface {
	Serialize() ([]byte, error)
	Create()
	Update()
	Delete()
}

func NewCardHub(mainHub *socket.MainHub) CardStorage {
	return &CardHub{
		mainHub: mainHub,
		mu:      sync.Mutex{},
		Cards:   make(map[string]*Card),
	}
}
