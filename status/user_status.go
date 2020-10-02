package status

// UserStates holds information about the connected Users
type UserStates map[string]User

// User holds information about the users status
type User struct {
	Username string
	Status   Status
}

// Status defines the status a user can set
type Status struct {
	Action   string
	Comments string
	Emotion  string
}
