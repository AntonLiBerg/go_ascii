package main

import (
	"fmt"
	aiapi "go_ascii/aiAPI"
	cmp "go_ascii/component"
)

func main() {
	world := NewWorld()
	ai := aiapi.New()
	aMap, entities, components := ai.GetAsciiMapAndEntitiesFromFile("./scenarios/demo/map.txt")

	for pos, ch := range aMap {
		eName := entities[ch]
		eComps := components[eName]
		world.AddEntity(pos, eComps)

	}
}

type World struct {
	nextEnt  int
	Entities []int
	Pos      map[int]cmp.Position
	Ascii    map[int]cmp.Ascii
}

func NewWorld() World {
	return World{
		nextEnt:  0,
		Entities: []int{},
		Pos:      map[int]cmp.Position{},
		Ascii:    map[int]cmp.Ascii{},
	}
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
			_, ok := compWithVals[cmp.C_POS]
			if !ok {
				return fmt.Errorf("missing expected components for %s: %v", name, vals)
			}
			w.Pos[eId] = cmp.Position{X: pos[0], Y: pos[1]}
		case cmp.C_ASCII:
			vals, ok := compWithVals[cmp.C_ASCII]
			if !ok {
				return fmt.Errorf("missing expected components for %s: %v", name, vals)
			}
			w.Ascii[eId] = cmp.Ascii{Ascii: []rune(vals[0])[0]}
		}

	}
	return nil
}
