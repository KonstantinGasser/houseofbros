package card

import (
	"encoding/json"
	"time"
)

type Task struct {
	UUID     string    `json:"uuid"`
	Title    string    `json:"title"`
	Author   string    `json:"author"`
	Assigned string    `json:"assigned"`
	Subject  string    `json:"subject"`
	IsDone   bool      `json:"is-done"`
	Iat      time.Time `json:"iat"`
}

func (u *Task) Serialize() ([]byte, error) {

	b, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return b, nil
}
