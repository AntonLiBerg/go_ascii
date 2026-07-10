package service

import (
	aiapi "go_ascii/aiAPI"
	wrld "go_ascii/world"
)

type ServiceDrawOnTerminal struct{}

func (s ServiceDrawOnTerminal) GetUpdateFunc(w wrld.World) UpdateFunc {
	return UpdateFunc{
		Order: 100,
		UpdateFunc: func(w *wrld.World) {
			if w.IterationNr == 1 || w.HasChanged {
				aiapi.UpdateTerminal(*w)
				w.HasChanged = false
			}
		},
	}
}
