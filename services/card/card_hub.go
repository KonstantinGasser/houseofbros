package card

import (
	"encoding/hex"
	"math/rand"
	"strings"
	"sync"
	"time"

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

func (hub *CardHub) UUID() (string, error) {

	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	b[7] = 0x40
	b[9] = 0x80

	_hex := hex.EncodeToString(b)
	hexString := strings.ReplaceAll(_hex, " ", "")

	uuid := []string{hexString[:8], "-", hexString[8:12], "-", hexString[12:16], "=", hexString[16:20], "=", hexString[20:]}
	return strings.Join(uuid, ""), nil
}

func (hub *CardHub) Serialize() ([]byte, error) {
	return nil, nil
}
