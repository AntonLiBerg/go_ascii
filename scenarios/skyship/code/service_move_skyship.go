package skyship

import (
	"go_ascii/internal/game"
	"go_ascii/internal/world"
)


type ServiceMoveSkyship struct{}

func (s ServiceMoveSkyship) GetUpdateFunc(w world.World) game.UpdateFunc {

	return game.UpdateFunc{
		Order: 1,
		UpdateFunc: func(w world.World) (world.World, error) {
			next := w
			shouldStart := next.ScenarioVariables["startship"] == "start"
			if !shouldStart {
				return next, nil
			}

			return next, nil
		},
	}
}

