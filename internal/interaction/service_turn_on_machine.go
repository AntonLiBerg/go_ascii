package interaction

import (
	"go_ascii/internal/game"
	wrld "go_ascii/internal/world"
	cmp "go_ascii/internal/world/component"
)

type ServiceTurnOnMachine struct{}

func (s ServiceTurnOnMachine) GetUpdateFunc(w wrld.World) game.UpdateFunc {
	interactKey := string(w.UserInputProfile.KeyInteract)
	if !w.UserInput[interactKey] {
		return game.UpdateFunc{Order: 1}
	}

	playerID := -1
	for _, eID := range w.Entities {
		if !w.HasComponent(eID, cmp.C_PLAYER) {
			continue
		}
		playerID = eID
		break
	}
	if playerID == -1 {
		return game.UpdateFunc{Order: 1}
	}

	neighbors := wrld.GetNeighbors(w, playerID, []cmp.ComponentName{cmp.C_MACHINE})
	machineID := -1
	if len(neighbors) == 0 {
		return game.UpdateFunc{}
	}
	if len(neighbors) == 1 {
		machineID = neighbors[0]
		return game.UpdateFunc{
			Order: 1,
			UpdateFunc: func(w *wrld.World) {
				machine := w.Machine[machineID]
				machine.IsOn = true
				w.Machine[machineID] = machine
				w.HasChanged = true
				w.UserInput[interactKey] = false
			},
		}
	}

	// Display choice menu
	return game.UpdateFunc{}
}
