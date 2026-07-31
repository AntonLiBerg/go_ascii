package world

type UserInputProfile struct {
	KeyQuitGame  string
	KeyExit      string
	KeyMoveUp    string
	KeyMoveLeft  string
	KeyMoveDown  string
	KeyMoveRight string
	KeyInteract  string
}

func NewUserInputProfile(values map[string]string) UserInputProfile {
	return UserInputProfile{
		KeyQuitGame:  values["quitgame"],
		KeyExit:      values["exit"],
		KeyMoveUp:    values["moveup"],
		KeyMoveLeft:  values["moveleft"],
		KeyMoveDown:  values["movedown"],
		KeyMoveRight: values["moveright"],
		KeyInteract:  values["interact"],
	}
}
