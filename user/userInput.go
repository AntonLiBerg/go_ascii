package user

type UserState int

const (
	S_playing UserState = iota
	S_quit
)

const ( 
	Key_quitGame = "quitgame"
	Key_moveDown = "movedown"
)

type UserInputProfile struct {
	KeyQuitGame string
	KeyMoveDown string
}

func NewUserInputProfileEmpty() UserInputProfile {
	return UserInputProfile{
		KeyQuitGame: "",
		KeyMoveDown: "",
	}
}

func NewUserInputProfile(userInputProfile map[string]string) UserInputProfile{
	return UserInputProfile{
		KeyQuitGame: userInputProfile[Key_quitGame],
		KeyMoveDown: userInputProfile[Key_moveDown],
	}
}
