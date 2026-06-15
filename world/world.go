package world

import (
	"fmt"
	cmp "go_ascii/component"
)

type World struct {
	userInput map[string]bool
	nextEnt   int
	Entities  []int
	Pos       map[int]cmp.Position
	Ascii     map[int]cmp.Ascii
}

func NewWorld() World {
	return World{
		userInput: map[string]bool{},
		nextEnt:   0,
		Entities:  []int{},
		Pos:       map[int]cmp.Position{},
		Ascii:     map[int]cmp.Ascii{},
	}
}
func (w *World) Clone() World {
	clone := World{
		userInput: make(map[string]bool, len(w.userInput)),
		nextEnt:   w.nextEnt,
		Entities:  append([]int(nil), w.Entities...),
		Pos:       make(map[int]cmp.Position, len(w.Pos)),
		Ascii:     make(map[int]cmp.Ascii, len(w.Ascii)),
	}

	for key, value := range w.userInput {
		clone.userInput[key] = value
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
	clear(w.userInput)
}

func (w *World) SetKeyDown(key string) {
	w.userInput[key] = true
}

func (w World) IsKeyDown(key string) bool {
	return w.userInput[key]
}

func (w *World) MakeNewEntityId() int {
	w.Entities = append(w.Entities, w.nextEnt)
	w.nextEnt++
	return w.nextEnt - 1
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
