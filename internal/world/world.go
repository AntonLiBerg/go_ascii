package world

import (
	"fmt"
	component "go_ascii/internal"
	"maps"
	"slices"
)

type World struct {
	UserInputProfile UserInputProfile
	KeyDown          string
	ShouldQuit       bool
	EByPos           map[component.Position]int
	HasChanged       bool
	IterationNr      int
	Entities         []int
	Pos              map[int]component.Position
	Ascii            map[int]component.Ascii
	Impassable       map[int]component.Impassable
	Player           map[int]component.Player
	Machine          map[int]component.Machine
}

func NewWorldEmpty() World {
	return World{
		EByPos:     make(map[component.Position]int),
		Pos:        make(map[int]component.Position),
		Ascii:      make(map[int]component.Ascii),
		Impassable: make(map[int]component.Impassable),
		Player:     make(map[int]component.Player),
		Machine:    make(map[int]component.Machine),
	}
}

func NewWorld(asciiMap map[[2]int]rune, entities map[rune]string, components map[string]map[string][]string) (World, error) {
	world := NewWorldEmpty()
	for pos, char := range asciiMap {
		if err := world.AddEntity(pos, components[entities[char]]); err != nil {
			return world, err
		}
	}
	return world, nil
}

func (w World) Clone() World {
	clone := w
	clone.Entities = slices.Clone(w.Entities)
	clone.EByPos = maps.Clone(w.EByPos)
	clone.Pos = maps.Clone(w.Pos)
	clone.Ascii = maps.Clone(w.Ascii)
	clone.Impassable = maps.Clone(w.Impassable)
	clone.Player = maps.Clone(w.Player)
	clone.Machine = maps.Clone(w.Machine)
	return clone
}

func (w *World) AddEntity(pos [2]int, components map[string][]string) error {
	eID := len(w.Entities)
	w.Entities = append(w.Entities, eID)

	for name, values := range components {
		switch name {
		case component.NamePosition:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			position := component.Position{X: pos[0], Y: pos[1]}
			w.Pos[eID] = position
			w.EByPos[position] = eID

		case component.NameASCII:
			if len(values) != 1 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			chars := []rune(values[0])
			if len(chars) != 1 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			w.Ascii[eID] = component.Ascii{Ascii: chars[0]}

		case component.NameImpassable:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			w.Impassable[eID] = component.Impassable{}

		case component.NamePlayer:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			w.Player[eID] = component.Player{}

		case component.NameMachine:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			w.Machine[eID] = component.Machine{}

		default:
			return fmt.Errorf("component does not exist %q", name)
		}
	}
	return nil
}
