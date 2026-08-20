package control

import (
	"fmt"
	component "go_ascii/internal"
	"go_ascii/internal/game"
	"go_ascii/internal/world"
	"slices"
	"strconv"
	"strings"
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
		return updateControlValue(w, w.ActiveControl.TargetEntityID, 1)
	case w.UserInputProfile.KeyMoveSelectPrev:
		return updateControlValue(w, w.ActiveControl.TargetEntityID, -1)
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
func updateControlValue(w world.World, entityID, delta int) game.UpdateFunc {
	if _, ok := w.ControlLabels[entityID]; ok {
		return updateControlLabelValue(entityID, delta)
	}
	if _, ok := w.ControlNumber[entityID]; ok {
		return updateControlNumberValue(entityID, delta)
	}
	return updateControlOptionsValue(entityID, delta)
}

func updateControlOptionsValue(entityID, delta int) game.UpdateFunc {
	return game.UpdateFunc{
		Order: 1,
		UpdateFunc: func(w world.World) (world.World, error) {
			next := w.Clone()
			controlOptions, ok := next.ControlOptions[entityID]
			if !ok {
				return next, fmt.Errorf("control options for entity %d not found", entityID)
			}
			if len(controlOptions.Options) == 0 {
				return next, fmt.Errorf("control options for entity %d are empty", entityID)
			}
			currentIndex := slices.Index(controlOptions.Options, controlOptions.Current)
			if currentIndex < 0 {
				return next, fmt.Errorf("current control option for entity %d not found", entityID)
			}

			nextIndex := (currentIndex + delta) % len(controlOptions.Options)
			if nextIndex < 0 {
				nextIndex += len(controlOptions.Options)
			}
			oldValue := controlOptions.Current
			controlOptions.Current = controlOptions.Options[nextIndex]
			next.ControlOptions[entityID] = controlOptions

			ascii, ok := next.Ascii[entityID]
			if !ok {
				return next, fmt.Errorf("ascii for control options entity %d not found", entityID)
			}
			ascii.Ascii = controlOptions.Current
			next.Ascii[entityID] = ascii
			next.HasChanged = oldValue != controlOptions.Current
			next.KeyDown = ""
			return next, nil
		},
	}
}
func updateControlLabelValue(entityID, delta int) game.UpdateFunc {
	return game.UpdateFunc{
		Order: 1,
		UpdateFunc: func(w world.World) (world.World, error) {
			next := w.Clone()
			label, ok := next.ControlLabels[entityID]
			if !ok {
				return next, fmt.Errorf("control label for entity %d not found", entityID)
			}
			if label.Operation != "append" {
				return next, fmt.Errorf("unsupported control label operation %q", label.Operation)
			}

			if delta < 0 {
				if len(label.History) == 0 {
					return next, nil
				}
				label.History = label.History[:len(label.History)-1]
			} else {
				parts := make([]string, len(label.Sources))
				for i, source := range label.Sources {
					if i < len(label.SourceEntityIDs) && label.SourceEntityIDs[i] >= 0 {
						value, err := controlValue(next, label.SourceEntityIDs[i])
						if err != nil {
							return next, err
						}
						parts[i] = value
					} else {
						parts[i] = source
					}
				}
				addition := strings.Join(parts, "")
				if len([]rune(label.Content))+len([]rune(addition)) > label.MaxLength {
					return next, nil
				}
				label.History = append(label.History, addition)
			}

			label.Content = strings.Join(label.History, "")
			if err := writeControlLabel(&next, label); err != nil {
				return next, err
			}
			next.ControlLabels[entityID] = label
			next.HasChanged = true
			next.KeyDown = ""
			return next, nil
		},
	}
}

func controlValue(w world.World, entityID int) (string, error) {
	if options, ok := w.ControlOptions[entityID]; ok {
		return string(options.Current), nil
	}
	if number, ok := w.ControlNumber[entityID]; ok {
		return strconv.Itoa(number.ValueCurrent), nil
	}
	return "", fmt.Errorf("entity %d has no control value", entityID)
}

func writeControlLabel(w *world.World, label component.ControlLabel) error {
	content := []rune(label.Content)
	if len(content) > len(label.EntityIDs) {
		content = content[:len(label.EntityIDs)]
	}
	for i, entityID := range label.EntityIDs {
		ascii, ok := w.Ascii[entityID]
		if !ok {
			return fmt.Errorf("control label entity %d has no ascii component", entityID)
		}
		ascii.Ascii = ' '
		if i < len(content) {
			ascii.Ascii = content[i]
		}
		w.Ascii[entityID] = ascii
	}
	return nil
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
			ascii.Ascii = rune('0' + controlNumber.ValueCurrent)
			next.Ascii[entityID] = ascii

			next.HasChanged = controlNumber.ValueCurrent != oldValue
			next.KeyDown = ""
			return next, nil
		},
	}
}
