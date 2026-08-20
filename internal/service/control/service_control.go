package control

import (
	"fmt"
	"go_ascii/internal/game"
	"go_ascii/internal/world"
)

type ServiceControl struct{}

func (s ServiceControl) GetUpdateFunc(w world.World) game.UpdateFunc {
	if w.UserInputProfile.ProfileType != world.ProfileTypeControl || w.KeyDown == "" {
		return game.UpdateFunc{}
	}
	if w.ActiveControl == nil {
		return game.UpdateFunc{}
	}

	if !w.EditingControl {
		switch w.KeyDown {
		case w.UserInputProfile.KeyMoveSelectNext:
			return game.UpdateFunc{
				Order: 1,
				UpdateFunc: func(w world.World) (world.World, error) {
					next, focused := w.WithNextControl()
					if focused {
						next.HasChanged = true
						next.KeyDown = ""
					}
					return next, nil
				},
			}
		case w.UserInputProfile.KeyMoveSelectPrev:
			return game.UpdateFunc{
				Order: 1,
				UpdateFunc: func(w world.World) (world.World, error) {
					next, focused := w.WithPreviousControl()
					if focused {
						next.HasChanged = true
						next.KeyDown = ""
					}
					return next, nil
				},
			}
		case w.UserInputProfile.KeySelect:
			return game.UpdateFunc{
				Order: 1,
				UpdateFunc: func(w world.World) (world.World, error) {
					next := w.Clone()
					next.EditingControl = true
					next.HasChanged = true
					next.KeyDown = ""
					return next, nil
				},
			}
		default:
			return game.UpdateFunc{}
		}
	}

	switch w.KeyDown {
	case w.UserInputProfile.KeyMoveSelectNext:
		return updateControlNumberValue(w.ActiveControl.TargetEntityID, 1)
	case w.UserInputProfile.KeyMoveSelectPrev:
		return updateControlNumberValue(w.ActiveControl.TargetEntityID, -1)
	case w.UserInputProfile.KeySelect:
		return game.UpdateFunc{
			Order: 1,
			UpdateFunc: func(w world.World) (world.World, error) {
				next := w.Clone()
				next.EditingControl = false
				next.HasChanged = true
				next.KeyDown = ""
				return next, nil
			},
		}
	default:
		return game.UpdateFunc{}
	}
}
func updateControlValue(w world.World,entityID int,delta int) game.UpdateFunc{
	if _,ok := w.ControlNumber[entityID]; ok{
		return updateControlNumberValue(entityID,delta)
	}
	return updateControlListValue(entityID,delta)

}
func updateControlListValue(entityID, delta int) game.UpdateFunc {
	return game.UpdateFunc{
		Order: 1,
		UpdateFunc: func(w world.World) (world.World, error) {
			next := w.Clone()
			controlList, ok := next.ControlList[entityID]
			if !ok {
				return next, fmt.Errorf("control list for entity %d not found", entityID)
			}
		}
	}
}
func updateControlNumberValue(entityID, delta int) game.UpdateFunc {
	return game.UpdateFunc{
		Order: 1,
		UpdateFunc: func(w world.World) (world.World, error) {
			next := w.Clone()
			controlNumber, ok := next.ControlNumber[entityID]
			if !ok {
				return next, fmt.Errorf("control number for entity %d not found", entityID)
			}
			oldValue := controlNumber.ValueCurrent
			if delta > 0 && controlNumber.ValueCurrent < controlNumber.ValueMax {
				controlNumber.ValueCurrent++
			} else if delta < 0 && controlNumber.ValueCurrent > controlNumber.ValueStart {
				controlNumber.ValueCurrent--
			}
			next.ControlNumber[entityID] = controlNumber

			ascii := next.Ascii[entityID]
			ascii.Ascii = rune('0'+controlNumber.ValueCurrent)
			next.Ascii[entityID] = ascii

			next.HasChanged = controlNumber.ValueCurrent != oldValue
			next.KeyDown = ""
			return next, nil
		},
	}
}
