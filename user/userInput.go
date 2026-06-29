package user

type UserState int

const (
	S_playing UserState = iota
	S_quit
)

const (
	Key_quitGame  = "quitgame"
	Key_moveUp    = "moveup"
	Key_moveLeft  = "moveleft"
	Key_moveDown  = "movedown"
	Key_moveRight = "moveright"
	KEY_INTERACT  = "interact"
)

type UserInputProfile struct {
	KeyQuitGame  string
	KeyMoveUp    string
	KeyMoveLeft  string
	KeyMoveDown  string
	KeyMoveRight string
	KeyInteract  string
}

func NewUserInputProfileEmpty() UserInputProfile {
	return UserInputProfile{
		KeyQuitGame:  "",
		KeyMoveUp:    "",
		KeyMoveLeft:  "",
		KeyMoveDown:  "",
		KeyMoveRight: "",
		KeyInteract:  "",
	}
}

func NewUserInputProfile(userInputProfile map[string]string) UserInputProfile {
	return UserInputProfile{
		KeyQuitGame:  userInputProfile[Key_quitGame],
		KeyMoveUp:    userInputProfile[Key_moveUp],
		KeyMoveLeft:  userInputProfile[Key_moveLeft],
		KeyMoveDown:  userInputProfile[Key_moveDown],
		KeyMoveRight: userInputProfile[Key_moveRight],
		KeyInteract:  userInputProfile[KEY_INTERACT],
	}
}
