package world

import (
	"fmt"
	cmp "go_ascii/component"
	usr "go_ascii/user"
)

type World struct {
	UserInputProfile usr.UserInputProfile
	StateUser        usr.UserState
	UserInput        map[string]bool
	NextEnt          int
	Entities         []int
	Pos              map[int]cmp.Position
	Ascii            map[int]cmp.Ascii
}

func NewWorldEmpty() World {
	return World{
		UserInputProfile: usr.NewUserInputProfileEmpty(),
		StateUser:        usr.S_playing,
		UserInput:        map[string]bool{},
		NextEnt:          0,
		Entities:         []int{},
		Pos:              map[int]cmp.Position{},
		Ascii:            map[int]cmp.Ascii{},
	}
}
func NewWorld(aMap map[[2]int]rune, entities map[rune]string, components map[string]map[cmp.ComponentName][]string) (World, error) {
	world := NewWorldEmpty()
	for pos, ch := range aMap {
		eName := entities[ch]
		eComps := components[eName]
		err := world.AddEntity(pos, eComps)
		if err != nil {
			return world, err
		}
	}
	return world, nil
}
func (w *World) Clone() World {
	clone := World{
		UserInput: make(map[string]bool, len(w.UserInput)),
		NextEnt:   w.NextEnt,
		Entities:  append([]int(nil), w.Entities...),
		Pos:       make(map[int]cmp.Position, len(w.Pos)),
		Ascii:     make(map[int]cmp.Ascii, len(w.Ascii)),
	}

	for key, value := range w.UserInput {
		clone.UserInput[key] = value
	}

	for id, pos := range w.Pos {
		clone.Pos[id] = pos
	}

	for id, ascii := range w.Ascii {
		clone.Ascii[id] = ascii
	}

	return clone
}
func (w *World) ClearUserInput() {
	clear(w.UserInput)
}

func (w *World) SetKeyDown(key string) {
	w.UserInput[key] = true
}

func (w World) IsKeyDown(key string) bool {
	return w.UserInput[key]
}

func (w *World) MakeNewEntityId() int {
	w.Entities = append(w.Entities, w.NextEnt)
	w.NextEnt++
	return w.NextEnt - 1
}

func (w *World) AddEntity(pos [2]int, compWithVals map[cmp.ComponentName][]string) error {
	eId := w.MakeNewEntityId()
	for name, vals := range compWithVals {
		switch name {
		case cmp.C_POS:
			w.Pos[eId] = cmp.Position{X: pos[0], Y: pos[1]}
		case cmp.C_ASCII:
			if len(vals) != 1 || len(vals[0]) != 1 {
				return fmt.Errorf("Required values are incorrect for %s", cmp.C_ASCII)
			}
			w.Ascii[eId] = cmp.Ascii{Ascii: []rune(vals[0])[0]}
		default:
			return fmt.Errorf("component does not exist %s", name)
		}
	}
	return nil
}
