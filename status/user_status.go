package status

import "log"

// Bro is something you cannot touch
type Bro struct {
	Uname    string        `json:"uname"`
	Activity string        `json:"activity"`
	Comment  string        `json:"comment"`
	Emotion  []interface{} `json:"emotion"`
}

// UpdateBro changes the state of a bro
func (bro *Bro) UpdateBro(activity, comment string, emotion []interface{}) *Bro {
	bro.Activity = activity
	bro.Comment = comment
	bro.Emotion = emotion
	log.Printf("[updated] bromotion")
	return bro
}

// NewBro welcome to the club
func NewBro(uname string) *Bro {
	return &Bro{
		Uname:    uname,
		Activity: "Being Awesome",
		Comment:  "this could be your ad",
		Emotion:  []interface{}{0},
	}
}
