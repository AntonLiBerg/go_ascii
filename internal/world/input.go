package world

type UserInputProfile struct {
	KeyQuitGame  string
	KeyMoveUp    string
	KeyMoveLeft  string
	KeyMoveDown  string
	KeyMoveRight string
	KeyInteract  string
}

func NewUserInputProfile(values map[string]string) UserInputProfile {
	return UserInputProfile{
		KeyQuitGame:  values["quitgame"],
		KeyMoveUp:    values["moveup"],
		KeyMoveLeft:  values["moveleft"],
		KeyMoveDown:  values["movedown"],
		KeyMoveRight: values["moveright"],
		KeyInteract:  values["interact"],
	}
}
