package service

import (
	ai "go_ascii/aiAPI"
	usr "go_ascii/user"
	wrld "go_ascii/world"
	cmp "go_ascii/component"
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
func (s ServiceMovePlayer) GetUpdateFunc(w wrld.World) UpdateFuncResult{
	return UpdateFuncResult{
		Order: 1,
		UpdateFunc: func(w *wrld.World){
			if w.UserInput[w.UserInputProfile.KeyMoveDown] {
			//Always only 1 entity!
				for k,_ := range w.EByTag[cmp.TAG_PLAYER]{
					TryGoToPosition(*w,k,cmp.Position{X:0, Y:+1})
				}
				w.UserInput[w.UserInputProfile.KeyMoveDown] = false
			}
		},
	}
}
func TryGoToPosition(w wrld.World, eMover int, posDelta cmp.Position) bool{

	moverPos := w.Pos[eMover]
	moverPos.X += posDelta.X
	moverPos.Y += posDelta.Y
	if !CanMakeMove(w,eMover,moverPos){
		return false
	}

	nPosEId := w.EByPos[moverPos]
	currentEntAtPos := w.Pos[nPosEId]
	currentEntAtPos.Y--

	w.Pos[eMover] = moverPos
	w.Pos[nPosEId] = currentEntAtPos
	return true
}
func CanMakeMove(w wrld.World, eMover int, posTarget cmp.Position) bool{
	_,err := w.Impassable[w.EByPos[posTarget]]
	if err {
		return false
	}
	return true
}

