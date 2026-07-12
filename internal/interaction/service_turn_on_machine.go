package interaction

import (
	"go_ascii/internal/game"
	wrld "go_ascii/internal/world"
	cmp "go_ascii/internal/world/component"
)

type ServiceTurnOnMachine struct{}

func (s ServiceTurnOnMachine) GetUpdateFunc(w wrld.World) game.UpdateFunc {
	interactKey := w.UserInputProfile.KeyInteract
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
				turnOnMachine(machineID, w)
			},
		}
	} else {
		return game.UpdateFunc{
			Order: 1,
			UpdateFunc: func(w *wrld.World) {
				for _, m := range neighbors {
					mch := w.Machine[m]
					pos := w.Pos[m]
					txt := pos.ToString() + ": " + string(mch.MachineType)
					cmenu := wrld.MenuChoice{
						Text: txt,
						Update: func(w *wrld.World) {
							turnOnMachine(machineID, w)
						},
					}
					w.MenuChoices.Choices = append(w.MenuChoices.Choices, cmenu)
				}
			},
		}
	}
}
func turnOnMachine(machineID int, w *wrld.World) {
	machine := w.Machine[machineID]
	machine.IsOn = true
	w.Machine[machineID] = machine
	w.HasChanged = true
	w.UserInput[w.UserInputProfile.KeyInteract] = false
}
