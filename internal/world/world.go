package world

import (
	"fmt"
	"maps"
	"slices"
)

const (
	componentPosition   = "pos"
	componentASCII      = "ascii"
	componentImpassable = "impassable"
	componentPlayer     = "player"
	componentMachine    = "machine"
)

type Position struct {
	X int
	Y int
}

type World struct {
	UserInputProfile UserInputProfile
	KeyDown          string
	ShouldQuit       bool
	EByPos           map[Position]int
	HasChanged       bool
	IterationNr      int
	Entities         []int
	Pos              map[int]Position
	Ascii            map[int]rune
	Impassable       map[int]struct{}
	Player           map[int]struct{}
	Machine          map[int]bool
}

func NewWorldEmpty() World {
	return World{
		EByPos:     make(map[Position]int),
		Pos:        make(map[int]Position),
		Ascii:      make(map[int]rune),
		Impassable: make(map[int]struct{}),
		Player:     make(map[int]struct{}),
		Machine:    make(map[int]bool),
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
		case componentPosition:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			position := Position{X: pos[0], Y: pos[1]}
			w.Pos[eID] = position
			w.EByPos[position] = eID

		case componentASCII:
			if len(values) != 1 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			chars := []rune(values[0])
			if len(chars) != 1 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			w.Ascii[eID] = chars[0]

		case componentImpassable:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			w.Impassable[eID] = struct{}{}

		case componentPlayer:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			w.Player[eID] = struct{}{}

		case componentMachine:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			w.Machine[eID] = false

		default:
			return fmt.Errorf("component does not exist %q", name)
		}
	}
	return nil
}
