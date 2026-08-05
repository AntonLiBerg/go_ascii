package control

import (
	"go_ascii/internal/game"
	"go_ascii/internal/world"
)

type ServiceControl struct{}

func (s ServiceControl) GetUpdateFunc(w world.World) game.UpdateFunc {

	if w.UserInputProfile.ProfileType != world.ProfileTypeControl || w.KeyDown == ""{
		return game.UpdateFunc{}
	}
	//
	// Is any control selected?
	//
	ctrlSelectedId, ok := w.GetSelectedControlId();
	if !ok {
		switch w.KeyDown{
		case w.UserInputProfile.KeyMoveSelectNext:
			return game.UpdateFunc{
				Order: 1,
				UpdateFunc: func(w world.World) (world.World, error) {
					nw := w.Clone();
					nw.FocusNextControl();
					return nw,nil;
				},
			}
		case w.UserInputProfile.KeyMoveSelectPrev:
			return game.UpdateFunc{
				Order: 1,
				UpdateFunc: func(w world.World) (world.World, error) {
					nw := w.Clone();
					nw.FocusPrevControl();
					return nw,nil;
				},
			}
		case w.UserInputProfile.KeySelect:
			return game.UpdateFunc{
				Order: 1,
				UpdateFunc: func(w world.World) (world.World, error) {
					nw := w.Clone();
					nw.SelectedControl = &w.FocusedControl;
					return nw,nil;
				},
			}
		default:
			return game.UpdateFunc{}
		}	
	}
  //
  // were there any changes done?
  //
  switch w.KeyDown{
  case w.UserInputProfile.KeyMoveSelectNext:
	  return game.UpdateFunc{
		  Order: 1,
		  UpdateFunc: func(w world.World) (world.World, error){
			  nw := w.Clone()
			  nrComp := nw.ControllNumber[ctrlSelectedId] 
			  nrComp.ValueCurrent++
			  return nw,nil;
		  },
	  }
  case w.UserInputProfile.KeyMoveSelectPrev:
	  return game.UpdateFunc{
		  Order: 1,
		  UpdateFunc: func(w world.World) (world.World, error){
			  nw := w.Clone()
			  nrComp := nw.ControllNumber[ctrlSelectedId] 
			  nrComp.ValueCurrent++
			  return nw,nil;
		  },
	  }
  default:
	  return game.UpdateFunc{}
  }
}
