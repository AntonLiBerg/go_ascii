package world

import (
	"fmt"
	component "go_ascii/internal"
	"go_ascii/internal/scenario"
	"maps"
	"slices"
)

type World struct {
	UserInputProfile   UserInputProfile
	InputProfiles      map[string]UserInputProfile
	InputProfileByRoom map[string]string
	KeyDown            string
	ShouldQuit         bool
	EByPos             map[component.Position]int
	Portals            map[component.Position]component.Position
	Terminals          map[component.Position]string
	ViewRoom           string
	HasChanged         bool
	IterationNr        int
	Entities           []int
	Pos                map[int]component.Position
	Layer              map[int]component.Layer
	Ascii              map[int]component.Ascii
	Impassable         map[int]component.Impassable
	Player             map[int]component.Player
	Door               map[int]component.Door
	Helm               map[int]component.Helm
	CommandTable       map[int]component.CommandTable
	BunkBed            map[int]component.BunkBed
	PrisonBars         map[int]component.PrisonBars
	Wall               map[int]component.Wall
	Window             map[int]component.Window
}

func NewWorldEmpty() World {
	return World{
		InputProfiles:      make(map[string]UserInputProfile),
		InputProfileByRoom: make(map[string]string),
		EByPos:             make(map[component.Position]int),
		Portals:            make(map[component.Position]component.Position),
		Terminals:          make(map[component.Position]string),
		Pos:                make(map[int]component.Position),
		Layer:              make(map[int]component.Layer),
		Ascii:              make(map[int]component.Ascii),
		Impassable:         make(map[int]component.Impassable),
		Player:             make(map[int]component.Player),
		Door:               make(map[int]component.Door),
		Helm:               make(map[int]component.Helm),
		CommandTable:       make(map[int]component.CommandTable),
		BunkBed:            make(map[int]component.BunkBed),
		PrisonBars:         make(map[int]component.PrisonBars),
		Wall:               make(map[int]component.Wall),
		Window:             make(map[int]component.Window),
	}
}

func NewWorld(asciiMap scenario.Map, entities map[rune]string, components map[string]map[string][]string, inputProfiles map[string]map[string]string) (World, error) {
	world := NewWorldEmpty()
	for name, values := range inputProfiles {
		world.InputProfiles[name] = NewUserInputProfile(values)
	}
	world.InputProfileByRoom = maps.Clone(asciiMap.InputProfiles)
	for roomName, profileName := range world.InputProfileByRoom {
		if _, ok := world.InputProfiles[profileName]; !ok {
			return world, fmt.Errorf("input profile %q for room %q does not exist", profileName, roomName)
		}
	}

	terminalRooms := make(map[string]struct{})
	for _, roomName := range asciiMap.Terminals {
		profileName := world.InputProfileByRoom[roomName]
		if world.InputProfiles[profileName].KeyExit == "" {
			return world, fmt.Errorf("terminal room %q input profile %q has no exit binding", roomName, profileName)
		}
		terminalRooms[roomName] = struct{}{}
	}

	roomNames := slices.Sorted(maps.Keys(asciiMap.Rooms))
	for _, roomName := range roomNames {
		room := asciiMap.Rooms[roomName]
		groundName := asciiMap.Ground[roomName]
		groundComponents, ok := components[groundName]
		if !ok {
			return world, fmt.Errorf("ground entity %q for room %q does not exist", groundName, roomName)
		}
		positions := slices.SortedFunc(maps.Keys(room), func(a, b [2]int) int {
			if a[1] != b[1] {
				return a[1] - b[1]
			}
			return a[0] - b[0]
		})
		for _, xy := range positions {
			position := component.Position{Room: roomName, X: xy[0], Y: xy[1]}
			if err := world.addEntityAtPosition(position, 0, groundComponents); err != nil {
				return world, err
			}

			char := room[xy]
			if char == ' ' {
				continue
			}
			if _, isTerminalRoom := terminalRooms[roomName]; isTerminalRoom {
				if err := world.addEntityAtPosition(position, 1, map[string][]string{
					component.NamePosition: {},
					component.NameASCII:    {string(char)},
				}); err != nil {
					return world, err
				}
				continue
			}
			if _, isPortal := asciiMap.Portals[position]; isPortal {
				continue
			}
			entityName, ok := entities[char]
			if !ok {
				return world, fmt.Errorf("map key %q in room %q has no entity", char, roomName)
			}
			if entityName == groundName {
				continue
			}
			if err := world.addEntityAtPosition(position, 1, components[entityName]); err != nil {
				return world, err
			}
		}
	}
	world.Portals = maps.Clone(asciiMap.Portals)
	world.Terminals = maps.Clone(asciiMap.Terminals)
	for position := range world.Terminals {
		isCommandTable := false
		for _, entityID := range GetEntitiesAtPosition(world, position) {
			if _, ok := world.CommandTable[entityID]; ok {
				isCommandTable = true
				break
			}
		}
		if !isCommandTable {
			return world, fmt.Errorf("terminal at %+v is not a commandtable entity", position)
		}
	}
	for playerID := range world.Player {
		world.SetInputProfileForRoom(world.Pos[playerID].Room)
		break
	}
	return world, nil
}

func (w *World) SetInputProfileForRoom(roomName string) bool {
	profileName, ok := w.InputProfileByRoom[roomName]
	if !ok {
		return false
	}
	profile, ok := w.InputProfiles[profileName]
	if !ok {
		return false
	}
	w.UserInputProfile = profile
	return true
}

func (w World) Clone() World {
	clone := w
	clone.Entities = slices.Clone(w.Entities)
	clone.InputProfiles = maps.Clone(w.InputProfiles)
	clone.InputProfileByRoom = maps.Clone(w.InputProfileByRoom)
	clone.EByPos = maps.Clone(w.EByPos)
	clone.Portals = maps.Clone(w.Portals)
	clone.Terminals = maps.Clone(w.Terminals)
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
	return w.addEntityAtPosition(component.Position{X: pos[0], Y: pos[1]}, layer, components)
}

func (w *World) AddEntityInRoom(room string, pos [2]int, components map[string][]string) error {
	return w.addEntityAtPosition(component.Position{Room: room, X: pos[0], Y: pos[1]}, 0, components)
}

func (w *World) addEntityAtPosition(position component.Position, layer int, components map[string][]string) error {
	eID := len(w.Entities)
	w.Entities = append(w.Entities, eID)
	w.Layer[eID] = component.Layer{Nr: layer}

	for name, values := range components {
		switch name {
		case component.NamePosition:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
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
