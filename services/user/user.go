package user

import "encoding/json"

type User interface {
	Serialize() ([]byte, error)
}

type StdUser struct {
	Username  string        `json:"username"`
	Action    string        `json:"action"`
	Note      string        `json:"note"`
	Emojies   []interface{} `json:"emojies"`
	Reactions []interface{} `json:"reactions"`
}

func (u *StdUser) Update(action, note string, emojies, reactions []interface{}) {
	u.Action = action
	u.Note = note
	u.Emojies = emojies
	u.Reactions = reactions
}

func (u *StdUser) Serialize() ([]byte, error) {

	b, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return b, nil
}
