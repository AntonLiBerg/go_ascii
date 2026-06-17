package service

import (
	ai "go_ascii/aiAPI"
	cmp "go_ascii/component"
	usr "go_ascii/user"
	wrld "go_ascii/world"
)

type UpdateFuncResult struct {
	Order      int
	UpdateFunc func(*wrld.World)
	Err        error
}
type IService interface {
	GetUpdateFunc(world wrld.World) UpdateFuncResult
}

type ServiceDrawOnTerminal struct{}

func (s ServiceDrawOnTerminal) GetUpdateFunc(w wrld.World) UpdateFuncResult {
	return UpdateFuncResult{
		Order: 100,
		UpdateFunc: func(w *wrld.World) {
			ai.UpdateTerminal(*w)
		},
	}
}

type ServiceQuitGame struct{}

func (s ServiceQuitGame) GetUpdateFunc(w wrld.World) UpdateFuncResult {
	return UpdateFuncResult{
		Order: 1,
		UpdateFunc: func(w *wrld.World) {
			if w.UserInput[w.UserInputProfile.KeyQuitGame] {
				w.StateUser = usr.S_quit
			}
		},
	}
}

type ServiceMovePlayer struct{}

func (s ServiceMovePlayer) GetUpdateFunc(w wrld.World) UpdateFuncResult {
	moveDelta := cmp.Position{}
	keyToClear := ""

	switch {
	case w.UserInput[w.UserInputProfile.KeyMoveUp]:
		moveDelta = cmp.Position{Y: -1}
		keyToClear = w.UserInputProfile.KeyMoveUp
	case w.UserInput[w.UserInputProfile.KeyMoveLeft]:
		moveDelta = cmp.Position{X: -1}
		keyToClear = w.UserInputProfile.KeyMoveLeft
	case w.UserInput[w.UserInputProfile.KeyMoveDown]:
		moveDelta = cmp.Position{Y: 1}
		keyToClear = w.UserInputProfile.KeyMoveDown
	case w.UserInput[w.UserInputProfile.KeyMoveRight]:
		moveDelta = cmp.Position{X: 1}
		keyToClear = w.UserInputProfile.KeyMoveRight
	default:
		return UpdateFuncResult{Order: 1}
	}

	return UpdateFuncResult{
		Order: 1,
		UpdateFunc: func(w *wrld.World) {
			for eID := range w.EByTag[cmp.TAG_PLAYER] {
				tryGoToPosition(w, eID, moveDelta)
			}
			w.UserInput[keyToClear] = false
		},
	}
}

func tryGoToPosition(w *wrld.World, eMover int, posDelta cmp.Position) bool {
	moverPos, ok := w.Pos[eMover]
	if !ok {
		return false
	}

	targetPos := cmp.Position{X: moverPos.X + posDelta.X, Y: moverPos.Y + posDelta.Y}
	targetID, ok := w.EByPos[targetPos]
	if !ok || !canMakeMove(w, targetID) {
		return false
	}

	w.Pos[eMover] = targetPos
	w.EByPos[targetPos] = eMover

	w.Pos[targetID] = moverPos
	w.EByPos[moverPos] = targetID
	return true
}

func canMakeMove(w *wrld.World, targetID int) bool {
	if _, blocked := w.Impassable[targetID]; blocked {
		return false
	}
	return true
}
