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

	ctrlSelectedID, ok := w.GetSelectedControlId()
	if !ok {
		switch w.KeyDown {
		case w.UserInputProfile.KeyMoveSelectNext:
			return game.UpdateFunc{
				Order: 1,
				UpdateFunc: func(w world.World) (world.World, error) {
					next := w.Clone()
					if next.FocusNextControl() {
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
					next := w.Clone()
					if next.FocusPrevControl() {
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
					next.SelectedControl = &next.FocusedControl
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
		return updateControlNumber(ctrlSelectedID, 1)
	case w.UserInputProfile.KeyMoveSelectPrev:
		return updateControlNumber(ctrlSelectedID, -1)
	default:
		return game.UpdateFunc{}
	}
}

func updateControlNumber(entityID, delta int) game.UpdateFunc {
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
			next.HasChanged = controlNumber.ValueCurrent != oldValue
			next.KeyDown = ""
			return next, nil
		},
	}
}
