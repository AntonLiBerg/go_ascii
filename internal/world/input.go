package world

const(
	ProfileTypeNone = "profiletypenone"
	ProfileTypeTerminal = "profiletypeterminal"
	ProfileTypeControl ="profiletypecontrol"
)
type UserInputProfile struct {
	ProfileType  string
	KeyQuitGame  string
	KeyExit      string
	KeyMoveUp    string
	KeyMoveLeft  string
	KeyMoveDown  string
	KeyMoveRight string
	KeyInteract  string
	KeyMoveSelectNext string
	KeyMoveSelectPrev string
	KeySelect string
}

func NewUserInputProfile(values map[string]string) UserInputProfile {
	return UserInputProfile{
		ProfileType: values["profiletype"],
		KeyQuitGame:  values["quitgame"],
		KeyExit:      values["exit"],
		KeyMoveUp:    values["moveup"],
		KeyMoveLeft:  values["moveleft"],
		KeyMoveDown:  values["movedown"],
		KeyMoveRight: values["moveright"],
		KeyInteract:  values["interact"],
		KeyMoveSelectNext: values["moveselectright"],
		KeyMoveSelectPrev: values["moveselectleft"],
		KeySelect: values["select"],
	}
}
