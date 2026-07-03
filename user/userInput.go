package user

type UserState int

const (
	S_playing UserState = iota
	S_INTERACT
	S_quit
)

type Interaction string

const (
	Key_quitGame  Interaction = "quitgame"
	Key_moveUp    Interaction = "moveup"
	Key_moveLeft  Interaction = "moveleft"
	Key_moveDown  Interaction = "movedown"
	Key_moveRight Interaction = "moveright"
	KEY_INTERACT  Interaction = "interact"
)

type UserInputProfile struct {
	KeyQuitGame  Interaction
	KeyMoveUp    Interaction
	KeyMoveLeft  Interaction
	KeyMoveDown  Interaction
	KeyMoveRight Interaction
	KeyInteract  Interaction
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
		KeyQuitGame:  Interaction(userInputProfile[string(Key_quitGame)]),
		KeyMoveUp:    Interaction(userInputProfile[string(Key_moveUp)]),
		KeyMoveLeft:  Interaction(userInputProfile[string(Key_moveLeft)]),
		KeyMoveDown:  Interaction(userInputProfile[string(Key_moveDown)]),
		KeyMoveRight: Interaction(userInputProfile[string(Key_moveRight)]),
		KeyInteract:  Interaction(userInputProfile[string(KEY_INTERACT)]),
	}
}
