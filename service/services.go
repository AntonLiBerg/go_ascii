package service

import (
	ai "go_ascii/aiAPI"
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

func (s *ServiceDrawOnTerminal) GetUpdateFunc(w wrld.World) UpdateFuncResult {
	return UpdateFuncResult{
		Order: 100,
		UpdateFunc: func(w *wrld.World) {

			ai.UpdateTerminal(*w)
		},
	}
}

type ServiceQuitGame struct{}

func (s *ServiceQuitGame) GetUpdateFunc(w wrld.World) UpdateFuncResult {
	return UpdateFuncResult{
		Order: 1,
		UpdateFunc: func(w *wrld.World) {
			if w.UserInput[w.UserInputProfile.KeyQuitGame] {
				w.StateUser = usr.S_quit
			}
		},
	}
}
