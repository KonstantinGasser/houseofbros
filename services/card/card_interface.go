package card

import (
	"sync"

	"github.com/KonstantinGasser/houseofbros/socket"
)

type CardStorage interface {
	Create()
	Update()
	Delete()
	UUID()
	Serialize() ([]byte, error)
}

func NewCardHub(mainHub *socket.MainHub) CardStorage {
	return &CardHub{
		mainHub: mainHub,
		mu:      sync.Mutex{},
		Cards:   make(map[string]*Card),
	}
}
