package quit

import (
	"go_ascii/internal/game"
	wrld "go_ascii/internal/world"
	usr "go_ascii/internal/world/user"
)

type ServiceQuitGame struct{}

func (s ServiceQuitGame) GetUpdateFunc(w wrld.World) game.UpdateFunc {
	return game.UpdateFunc{
		Order: 1,
		UpdateFunc: func(w *wrld.World) {
			if w.UserInput[string(w.UserInputProfile.KeyQuitGame)] {
				w.SetUserState(usr.S_quit, true)
			}
		},
	}
}
