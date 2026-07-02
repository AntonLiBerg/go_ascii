package user

type UserState int

const (
	S_playing UserState = iota
	S_INTERACT
	S_quit
)

type Interaction string
const (
	Key_quitGame Interaction = "quitgame"
	Key_moveUp  Interaction   = "moveup"
	Key_moveLeft  Interaction = "moveleft"
	Key_moveDown  Interaction = "movedown"
	Key_moveRight  Interaction= "moveright"
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
		KeyQuitGame:  Key_quitGame,
		KeyMoveUp:    Key_moveUp,
		KeyMoveLeft:  Key_moveLeft,
		KeyMoveDown:  Key_moveDown,
		KeyMoveRight: Key_moveRight,
		KeyInteract:  KEY_INTERACT,
	}
}
