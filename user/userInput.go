package user

type UserState int

const (
	S_playing UserState = iota
	S_quit
)

type UserInputProfile struct {
	KeyQuitGame string
}

func NewUserInputProfileEmpty() UserInputProfile {
	return UserInputProfile{
		KeyQuitGame: "",
	}
}
