package world

import (
	"fmt"
	component "go_ascii/internal"
	"go_ascii/internal/scenario"
	"maps"
	"slices"
	"strconv"
)

type ControlNode struct {
	SelectableEntityID int
	TargetEntityID     int
	Next               *ControlNode
	Prev               *ControlNode
}

type World struct {
	Room                   string
	UserInputProfile       UserInputProfile
	InputProfiles          map[string]UserInputProfile
	InputProfileByRoom     map[string]string
	KeyDown                string
	InfoboxContent         string
	ShouldQuit             bool
	ActiveControl          *ControlNode
	EditingControl         bool
	ControlLists           map[string]*ControlNode
	EByPos                 map[component.Position]int
	Portals                map[component.Position]component.Position
	Terminals              map[component.Position]string
	UILayout               []string
	UIs                    []string
	UIContent              map[string][]string
	UIEmpty                map[string][]string
	InfoboxRoomBaseContent map[string]string
	HasChanged             bool
	IterationNr            int
	Entities               []int
	Pos                    map[int]component.Position
	Layer                  map[int]component.Layer
	Ascii                  map[int]component.Ascii
	Impassable             map[int]component.Impassable
	Player                 map[int]component.Player
	Interactable           map[int]component.Interactable
	ControlNumber          map[int]component.ControlNumber
	Selectable             map[int]component.Selectable
}

func NewWorldEmpty() World {
	return World{
		InputProfiles:          make(map[string]UserInputProfile),
		InputProfileByRoom:     make(map[string]string),
		EByPos:                 make(map[component.Position]int),
		Portals:                make(map[component.Position]component.Position),
		Terminals:              make(map[component.Position]string),
		UIContent:              make(map[string][]string),
		UIEmpty:                make(map[string][]string),
		InfoboxRoomBaseContent: make(map[string]string),
		Pos:                    make(map[int]component.Position),
		Layer:                  make(map[int]component.Layer),
		Ascii:                  make(map[int]component.Ascii),
		Impassable:             make(map[int]component.Impassable),
		Player:                 make(map[int]component.Player),
		Interactable:           make(map[int]component.Interactable),
		ControlNumber:          make(map[int]component.ControlNumber),
		Selectable:             make(map[int]component.Selectable),
		ControlLists:           make(map[string]*ControlNode),
	}
}

func NewWorld(asciiMap scenario.FileMap, entities map[rune]string, components map[string]map[string][]string, inputProfiles map[string]map[string]string) (World, error) {
	return newWorld(asciiMap, entities, components, inputProfiles, nil, nil)
}

func NewWorldWithUI(asciiMap scenario.FileMap, entities map[rune]string, components map[string]map[string][]string, inputProfiles map[string]map[string]string, uiLayout []string, uis []string) (World, error) {
	return newWorld(asciiMap, entities, components, inputProfiles, uiLayout, uis)
}

func newWorld(asciiMap scenario.FileMap, entities map[rune]string, components map[string]map[string][]string, inputProfiles map[string]map[string]string, uiLayout []string, uis []string) (World, error) {
	world := NewWorldEmpty()
	world.UILayout = slices.Clone(uiLayout)
	world.UIs = slices.Clone(uis)
	for name, lines := range asciiMap.UIContent {
		world.UIContent[name] = slices.Clone(lines)
		world.UIEmpty[name] = slices.Clone(lines)
	}
	world.InfoboxRoomBaseContent = maps.Clone(asciiMap.InfoboxRoomBaseContent)
	for name, values := range inputProfiles {
		world.InputProfiles[name] = NewUserInputProfile(values)
	}
	roomEntityIDs := make(map[string]map[rune]int)
	entityIDsByName := make(map[string]int)
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

	if len(asciiMap.Rooms) == 0 {
		return world, fmt.Errorf("world has no rooms")
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
				entityName, hasEntity := entities[char]
				if allowed, scoped := asciiMap.EntityGroups[roomName]; scoped {
					_, explicitlyAllowed := allowed[char]
					hasEntity = hasEntity && explicitlyAllowed
				}
				entityID := len(world.Entities)
				if hasEntity {
					entityComponents, exists := components[entityName]
					if !exists {
						return world, fmt.Errorf("entity %q for map key %q does not have a definition", entityName, char)
					}
					if err := world.addEntityAtPosition(position, 1, entityComponents); err != nil {
						return world, err
					}
					entityIDsByName[entityName] = entityID
				} else if err := world.addEntityAtPosition(position, 1, map[string][]string{
					component.NamePosition: {},
					component.NameASCII:    {string(char)},
				}); err != nil {
					return world, err
				}
				if hasEntity {
					if roomEntityIDs[roomName] == nil {
						roomEntityIDs[roomName] = make(map[rune]int)
					}
					roomEntityIDs[roomName][char] = entityID
				}
				continue
			}
			if _, isPortal := asciiMap.Portals[position]; isPortal {
				continue
			}
			if allowed, scoped := asciiMap.EntityGroups[roomName]; scoped {
				if _, ok := allowed[char]; !ok {
					return world, fmt.Errorf("map key %q in room %q is not in its entity groups", char, roomName)
				}
			}
			entityName, ok := entities[char]
			if !ok {
				return world, fmt.Errorf("map key %q in room %q has no entity", char, roomName)
			}
			entityComponents, ok := components[entityName]
			if !ok {
				return world, fmt.Errorf("entity %q for map key %q does not have a definition", entityName, char)
			}
			if entityName == groundName {
				continue
			}
			entityID := len(world.Entities)
			if err := world.addEntityAtPosition(position, 1, entityComponents); err != nil {
				return world, err
			}
			entityIDsByName[entityName] = entityID
			if roomEntityIDs[roomName] == nil {
				roomEntityIDs[roomName] = make(map[rune]int)
			}
			roomEntityIDs[roomName][char] = entityID
		}
	}
	for entityID, selectable := range world.Selectable {
		targetID, ok := entityIDsByName[selectable.TargetEntityName]
		if !ok {
			return world, fmt.Errorf("selectable entity %d targets unknown entity %q", entityID, selectable.TargetEntityName)
		}
		selectable.TargetEntityId = targetID
		selectable.TargetEntityName = ""
		world.Selectable[entityID] = selectable
	}
	for roomName, markers := range asciiMap.SelectableOrder {
		ids := make([]int, 0, len(markers))
		for _, marker := range markers {
			entityID, ok := roomEntityIDs[roomName][marker]
			if !ok {
				return world, fmt.Errorf("selectable order for room %q references missing entity %q", roomName, marker)
			}
			if _, ok := world.Selectable[entityID]; !ok {
				return world, fmt.Errorf("entity %q in room %q is not selectable", marker, roomName)
			}
			ids = append(ids, entityID)
		}
		if len(ids) > 0 {
			world.ControlLists[roomName] = circularControlNodes(ids, world.Selectable)
		}
	}
	world.Portals = maps.Clone(asciiMap.Portals)
	world.Terminals = maps.Clone(asciiMap.Terminals)
	for position := range world.Terminals {
		isTerminalInteraction := false
		for _, entityID := range GetEntitiesAtPosition(world, position) {
			interactable, ok := world.Interactable[entityID]
			if ok && interactable.InteractionType == component.InteractionTypeTerminal {
				isTerminalInteraction = true
				break
			}
		}
		if !isTerminalInteraction {
			return world, fmt.Errorf("terminal at %+v does not have a terminal interaction", position)
		}
	}
	// The active room follows the single player after all entities are built.
	if len(world.Player) != 1 {
		return world, fmt.Errorf("world must have exactly one player, got %d", len(world.Player))
	}
	playerID := GetPlayerID(world)
	if position, ok := world.Pos[playerID]; ok {
		world.Room = position.Room
		world.SetInputProfileForRoom(position.Room)
		world.InfoboxContent = world.InfoboxRoomBaseContent[world.Room]
	}
	return world, nil
}
func (w *World) FocusNextControl() bool {
	if w.ActiveControl == nil || w.ActiveControl.Next == nil {
		return false
	}
	w.ActiveControl = w.ActiveControl.Next
	return true
}

func (w *World) FocusPrevControl() bool {
	if w.ActiveControl == nil || w.ActiveControl.Prev == nil {
		return false
	}
	w.ActiveControl = w.ActiveControl.Prev
	return true
}
func (w *World) SetInputProfileForRoom(roomName string) bool {
	profile, ok := w.InputProfileForRoom(roomName)
	if !ok {
		return false
	}
	w.UserInputProfile = profile
	w.EditingControl = false
	w.ActiveControl = w.ControlLists[roomName]
	return true
}

func (w World) InputProfileForRoom(roomName string) (UserInputProfile, bool) {
	profileName, ok := w.InputProfileByRoom[roomName]
	if !ok {
		return UserInputProfile{}, false
	}
	profile, ok := w.InputProfiles[profileName]
	if !ok {
		return UserInputProfile{}, false
	}
	return profile, true
}

func (w World) Clone() World {
	clone := w
	clone.Entities = slices.Clone(w.Entities)
	clone.InputProfiles = maps.Clone(w.InputProfiles)
	clone.InputProfileByRoom = maps.Clone(w.InputProfileByRoom)
	clone.EByPos = maps.Clone(w.EByPos)
	clone.Portals = maps.Clone(w.Portals)
	clone.Terminals = maps.Clone(w.Terminals)
	clone.UILayout = slices.Clone(w.UILayout)
	clone.UIs = slices.Clone(w.UIs)
	clone.UIContent = make(map[string][]string, len(w.UIContent))
	for name, lines := range w.UIContent {
		clone.UIContent[name] = slices.Clone(lines)
	}
	clone.UIEmpty = make(map[string][]string, len(w.UIEmpty))
	for name, lines := range w.UIEmpty {
		clone.UIEmpty[name] = slices.Clone(lines)
	}
	clone.InfoboxRoomBaseContent = maps.Clone(w.InfoboxRoomBaseContent)
	clone.Pos = maps.Clone(w.Pos)
	clone.Layer = maps.Clone(w.Layer)
	clone.Ascii = maps.Clone(w.Ascii)
	clone.Impassable = maps.Clone(w.Impassable)
	clone.Player = maps.Clone(w.Player)
	clone.Interactable = maps.Clone(w.Interactable)
	clone.ControlNumber = maps.Clone(w.ControlNumber)
	clone.Selectable = maps.Clone(w.Selectable)
	clone.ControlLists = maps.Clone(w.ControlLists)
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
	var ascii component.Ascii
	var hasASCII bool
	var hasImpassable, hasPlayer bool
	var interactable component.Interactable
	var hasInteractable bool
	var controlNumber component.ControlNumber
	var hasControlNumber bool
	var selectable component.Selectable
	var hasSelectable bool

	for name, values := range components {
		switch name {
		case component.NamePosition:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}

		case component.NameASCII:
			if len(values) != 1 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			chars := []rune(values[0])
			if len(chars) != 1 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			ascii = component.Ascii{Ascii: chars[0]}
			hasASCII = true

		case component.NameImpassable:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			hasImpassable = true

		case component.NamePlayer:
			if len(values) != 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			hasPlayer = true

		case component.NameInteractable:
			if len(values) != 1 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			interactable = component.Interactable{InteractionType: values[0]}
			hasInteractable = true

		case component.NameControlTypeNumber:
			if len(values) != 2 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			valueStart, err := strconv.Atoi(values[0])
			if err != nil {
				return fmt.Errorf("invalid values for component %q", name)
			}
			valueMax, err := strconv.Atoi(values[1])
			if err != nil || valueStart > valueMax {
				return fmt.Errorf("invalid values for component %q", name)
			}
			controlNumber = component.ControlNumber{
				ValueStart:   valueStart,
				ValueCurrent: valueStart,
				ValueMax:     valueMax,
			}
			hasControlNumber = true

		case component.NameSelectable:
			if len(values) != 4 || values[3] == "" {
				return fmt.Errorf("invalid values for component %q", name)
			}
			unfocused := []rune(values[0])
			focused := []rune(values[1])
			selected := []rune(values[2])
			if len(unfocused) != 1 || len(focused) != 1 || len(selected) != 1 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			selectable = component.Selectable{
				UnfocusedASCII:   unfocused[0],
				FocusedASCII:     focused[0],
				SelectedASCII:    selected[0],
				TargetEntityName: values[3],
			}
			hasSelectable = true

		default:
			return fmt.Errorf("component does not exist %q", name)
		}
	}

	w.Entities = append(w.Entities, eID)
	w.Layer[eID] = component.Layer{Nr: layer}
	w.Pos[eID] = position
	w.EByPos[position] = eID
	if hasASCII {
		w.Ascii[eID] = ascii
	}
	if hasImpassable {
		w.Impassable[eID] = component.Impassable{}
	}
	if hasPlayer {
		w.Player[eID] = component.Player{}
	}
	if hasInteractable {
		w.Interactable[eID] = interactable
	}
	if hasControlNumber {
		w.ControlNumber[eID] = controlNumber
	}
	if hasSelectable {
		w.Selectable[eID] = selectable
	}
	return nil
}

func circularControlNodes(entityIDs []int, selectables map[int]component.Selectable) *ControlNode {
	nodes := make([]*ControlNode, len(entityIDs))
	for i, entityID := range entityIDs {
		nodes[i] = &ControlNode{
			SelectableEntityID: entityID,
			TargetEntityID:     selectables[entityID].TargetEntityId,
		}
	}
	for i, node := range nodes {
		node.Next = nodes[(i+1)%len(nodes)]
		node.Prev = nodes[(i+len(nodes)-1)%len(nodes)]
	}
	return nodes[0]
}
