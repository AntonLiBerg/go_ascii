package quit

import (
	"go_ascii/internal/game"
	"go_ascii/internal/world"
)

type ServiceQuitGame struct{}

func (s ServiceQuitGame) GetUpdateFunc(w world.World) game.UpdateFunc {
	quitKey := w.UserInputProfile.KeyQuitGame
	if quitKey == "" || w.KeyDown != quitKey {
		return game.UpdateFunc{Order: 1}
	}

	return game.UpdateFunc{
		Order: 1,
		UpdateFunc: func(w world.World) (world.World, error) {
			next := w
			next.ShouldQuit = true
			return next, nil
		},
	}
}
