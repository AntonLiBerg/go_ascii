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
	Layer            map[int]component.Layer
	Ascii            map[int]component.Ascii
	Impassable       map[int]component.Impassable
	Player           map[int]component.Player
	Door             map[int]component.Door
	Helm             map[int]component.Helm
	CommandTable     map[int]component.CommandTable
	BunkBed          map[int]component.BunkBed
	PrisonBars       map[int]component.PrisonBars
	Wall             map[int]component.Wall
	Window           map[int]component.Window
}

func NewWorldEmpty() World {
	return World{
		EByPos:       make(map[component.Position]int),
		Pos:          make(map[int]component.Position),
		Layer:        make(map[int]component.Layer),
		Ascii:        make(map[int]component.Ascii),
		Impassable:   make(map[int]component.Impassable),
		Player:       make(map[int]component.Player),
		Door:         make(map[int]component.Door),
		Helm:         make(map[int]component.Helm),
		CommandTable: make(map[int]component.CommandTable),
		BunkBed:      make(map[int]component.BunkBed),
		PrisonBars:   make(map[int]component.PrisonBars),
		Wall:         make(map[int]component.Wall),
		Window:       make(map[int]component.Window),
	}
}

func NewWorld(asciiMap map[int]map[[2]int]rune, entities map[rune]string, components map[string]map[string][]string) (World, error) {
	world := NewWorldEmpty()
	layers := slices.Sorted(maps.Keys(asciiMap))
	for _, layer := range layers {
		for pos, char := range asciiMap[layer] {
			entityName, ok := entities[char]
			if !ok {
				return world, fmt.Errorf("map key %q on layer %d has no entity", char, layer)
			}
			if err := world.AddEntityAtLayer(pos, layer, components[entityName]); err != nil {
				return world, err
			}
		}
	}
	return world, nil
}

func (w World) Clone() World {
	clone := w
	clone.Entities = slices.Clone(w.Entities)
	clone.EByPos = maps.Clone(w.EByPos)
	clone.Pos = maps.Clone(w.Pos)
	clone.Layer = maps.Clone(w.Layer)
	clone.Ascii = maps.Clone(w.Ascii)
	clone.Impassable = maps.Clone(w.Impassable)
	clone.Player = maps.Clone(w.Player)
	clone.Door = maps.Clone(w.Door)
	clone.Helm = maps.Clone(w.Helm)
	clone.CommandTable = maps.Clone(w.CommandTable)
	clone.BunkBed = maps.Clone(w.BunkBed)
	clone.PrisonBars = maps.Clone(w.PrisonBars)
	clone.Wall = maps.Clone(w.Wall)
	clone.Window = maps.Clone(w.Window)
	return clone
}

func (w *World) AddEntity(pos [2]int, components map[string][]string) error {
	return w.AddEntityAtLayer(pos, 0, components)
}

func (w *World) AddEntityAtLayer(pos [2]int, layer int, components map[string][]string) error {
	eID := len(w.Entities)
	w.Entities = append(w.Entities, eID)
	w.Layer[eID] = component.Layer{Nr: layer}

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

		case component.NameDoor:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			w.Door[eID] = component.Door{IsInteractable: true}

		case component.NameHelm:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			w.Helm[eID] = component.Helm{IsInteractable: true}

		case component.NameCommandTable:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			w.CommandTable[eID] = component.CommandTable{IsInteractable: true}

		case component.NameBunkBed:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			w.BunkBed[eID] = component.BunkBed{}

		case component.NamePrisonBars:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			w.PrisonBars[eID] = component.PrisonBars{}

		case component.NameWall:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			w.Wall[eID] = component.Wall{}

		case component.NameWindow:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			w.Window[eID] = component.Window{}

		default:
			return fmt.Errorf("component does not exist %q", name)
		}
	}
	return nil
}
