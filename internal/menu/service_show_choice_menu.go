package menu

import (
	"go_ascii/internal/game"
	wrld "go_ascii/internal/world"
)

type ServiceShowChoiceMenu struct{}

func (s ServiceShowChoiceMenu) GetUpdateFunc(w wrld.World) game.UpdateFunc {
	if !w.MenuChoices.ShouldShow || w.MenuChoices.IsOpen {
		return game.UpdateFunc{Order: 1}
	}

	return game.UpdateFunc{
		Order: 1,
		UpdateFunc: func(w *wrld.World) {
			w.MenuChoices.IsOpen = true
			w.HasChanged = true
		},
	}
}
