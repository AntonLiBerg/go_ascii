package interaction

import (
	"go_ascii/internal/game"
	"go_ascii/internal/world"
)

type ServiceTurnOnMachine struct{}

func (s ServiceTurnOnMachine) GetUpdateFunc(w world.World) game.UpdateFunc {
	interactKey := w.UserInputProfile.KeyInteract
	if interactKey == "" || w.KeyDown != interactKey {
		return game.UpdateFunc{Order: 1}
	}

	playerID := -1
	for eID := range w.Player {
		playerID = eID
		break
	}

	machineID := -1
	if playerID != -1 {
		for _, neighborID := range world.GetNeighbors(w, playerID) {
			if _, isMachine := w.Machine[neighborID]; isMachine {
				machineID = neighborID
				break
			}
		}
	}

	return game.UpdateFunc{
		Order: 1,
		UpdateFunc: func(w *world.World) {
			if machineID != -1 {
				machine := w.Machine[machineID]
				machine.IsOn = true
				w.Machine[machineID] = machine
				w.HasChanged = true
			}
			if w.KeyDown == interactKey {
				w.KeyDown = ""
			}
		},
	}
}
