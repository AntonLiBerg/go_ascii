package service

import (
	aiapi "go_ascii/aiAPI"
	cmp "go_ascii/component"
	wrld "go_ascii/world"
)

type ServiceTurnOnMachine struct{}

func (s ServiceTurnOnMachine) GetUpdateFunc(w wrld.World) UpdateFunc {
	interactKey := string(w.UserInputProfile.KeyInteract)
	if !w.UserInput[interactKey] {
		return UpdateFunc{Order: 1}
	}

	playerID := -1
	for eID := range w.EByTag[cmp.TAG_PLAYER] {
		playerID = eID
		break
	}
	if playerID == -1 {
		return UpdateFunc{Order: 1}
	}

	neighbors := aiapi.GetNeighbors(w, playerID, []cmp.ComponentName{cmp.C_MACHINE})
	machineID := -1
	if len(neighbors) == 0 {
		return UpdateFunc{}
	}
	if len(neighbors) == 1 {
		machineID = neighbors[0]
		return UpdateFunc{
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
	return UpdateFunc{}
}
