package card

import (
	"log"
	"sync"

	"github.com/KonstantinGasser/houseofbros/socket"
)

type CardStorage interface {
	Create(v map[string]interface{})
	Update(uuid string) ([]byte, error)
	Delete()
	GenerateUUID() (string, error)
	Serialize() ([]byte, error)
}

func NewCardHub(mainHub *socket.MainHub) CardStorage {
	log.Printf("[created] new CardStorage\n")
	return &CardHub{
		mainHub: mainHub,
		mu:      sync.Mutex{},
		Cards:   make(map[string]Card),
	}
}
