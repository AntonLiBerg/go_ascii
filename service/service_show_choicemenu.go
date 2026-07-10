package service

import (
	wrld "go_ascii/world"
)

type ServiceShowChoiceMenu struct{}

func (s ServiceShowChoiceMenu) GetUpdateFunc(w wrld.World) UpdateFunc {
	if !w.MenuChoices.ShouldShow || w.MenuChoices.IsOpen {
		return UpdateFunc{Order: 1}
	}

	return UpdateFunc{
		Order: 1,
		UpdateFunc: func(w *wrld.World) {
			w.MenuChoices.IsOpen = true
			w.HasChanged = true
		},
	}
}
