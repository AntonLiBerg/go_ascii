package service

import (
	usr "go_ascii/user"
	wrld "go_ascii/world"
)

type ServiceQuitGame struct{}

func (s ServiceQuitGame) GetUpdateFunc(w wrld.World) UpdateFunc {
	return UpdateFunc{
		Order: 1,
		UpdateFunc: func(w *wrld.World) {
			if w.UserInput[string(w.UserInputProfile.KeyQuitGame)] {
				w.SetUserState(usr.S_quit, true)
			}
		},
	}
}
