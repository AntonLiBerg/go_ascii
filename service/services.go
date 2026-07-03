package service

import (
	"fmt"
	aiapi "go_ascii/aiAPI"
	cmp "go_ascii/component"
	usr "go_ascii/user"
	wrld "go_ascii/world"
)

type UpdateFunc struct {
	Order      int
	UpdateFunc func(*wrld.World)
	Err        error
}
type IService interface {
	GetUpdateFunc(world wrld.World) UpdateFunc
}

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

type ServiceQuitGame struct{}

func (s ServiceQuitGame) GetUpdateFunc(w wrld.World) UpdateFunc {
	return UpdateFunc{
		Order: 1,
		UpdateFunc: func(w *wrld.World) {
			if w.UserInput[string(w.UserInputProfile.KeyQuitGame)] {
				w.StateUser = usr.S_quit
			}
		},
	}
}

type ServiceMovePlayer struct{}

func (s ServiceMovePlayer) GetUpdateFunc(w wrld.World) UpdateFunc {
	moveDelta := cmp.Position{}
	keyToClear := ""

	switch {
	case w.UserInput[string(w.UserInputProfile.KeyMoveUp)]:
		moveDelta = cmp.Position{Y: -1}
		keyToClear = string(w.UserInputProfile.KeyMoveUp)
	case w.UserInput[string(w.UserInputProfile.KeyMoveLeft)]:
		moveDelta = cmp.Position{X: -1}
		keyToClear = string(w.UserInputProfile.KeyMoveLeft)
	case w.UserInput[string(w.UserInputProfile.KeyMoveDown)]:
		moveDelta = cmp.Position{Y: 1}
		keyToClear = string(w.UserInputProfile.KeyMoveDown)
	case w.UserInput[string(w.UserInputProfile.KeyMoveRight)]:
		moveDelta = cmp.Position{X: 1}
		keyToClear = string(w.UserInputProfile.KeyMoveRight)
	default:
		return UpdateFunc{Order: 1}
	}

	return UpdateFunc{
		Order: 1,
		UpdateFunc: func(w *wrld.World) {
			w.HasChanged = true
			for eID := range w.EByTag[cmp.TAG_PLAYER] {
				if err := tryGoToPosition(w, eID, moveDelta); err != nil {
					panic(err)
				}
			}
			w.UserInput[keyToClear] = false
		},
	}
}

func tryGoToPosition(w *wrld.World, eMover int, posDelta cmp.Position) error {
	moverPos, ok := w.Pos[eMover]
	if !ok {
		return fmt.Errorf("Mover entity not found!")
	}

	targetPos := cmp.Position{X: moverPos.X + posDelta.X, Y: moverPos.Y + posDelta.Y}
	targetID, ok := w.EByPos[targetPos]
	if !ok {
		return nil
	}
	if !canMakeMove(w, targetID) {
		return nil
	}

	w.Pos[eMover] = targetPos
	w.EByPos[targetPos] = eMover

	w.Pos[targetID] = moverPos
	w.EByPos[moverPos] = targetID
	return nil
}

func canMakeMove(w *wrld.World, targetID int) bool {
	if _, blocked := w.Impassable[targetID]; blocked {
		return false
	}
	return true
}

type ServiceTurnOnMachine struct{}

func (s ServiceTurnOnMachine) GetUpdateFunc(w wrld.World) UpdateFunc {
	if !w.UserInput[string(w.UserInputProfile.KeyInteract)] {
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
	if len(neighbors) == 1 {
		machineID = neighbors[0]
	}

	return UpdateFunc{
		Order: 1,
		UpdateFunc: func(w *wrld.World) {
			w.UserInput[string(w.UserInputProfile.KeyInteract)] = false
			if machineID == -1 {
				return
			}

			machine := w.Machine[machineID]
			if machine.IsOn {
				return
			}

			machine.IsOn = true
			w.Machine[machineID] = machine
			w.HasChanged = true
		},
	}
}
