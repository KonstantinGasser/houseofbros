package socket

import "encoding/json"

type Event interface {
	Serialize() ([]byte, error)
}

type EventUser struct {
	Type string `json:"event"`
	User []byte `json:"user"`
}

func (evt EventUser) Serialize() ([]byte, error) {
	b, err := json.Marshal(evt)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type EventReaction struct {
	Type     string      `json:"event"`
	From     string      `json:"from-user"`
	To       string      `json:"to-user"`
	Reaction interface{} `json:"reaction"`
}

func (evt EventReaction) Serialize() ([]byte, error) {
	b, err := json.Marshal(evt)
	if err != nil {
		return nil, err
	}
	return b, nil
}
