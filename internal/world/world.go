package world

import (
	"fmt"
	component "go_ascii/internal"
	"go_ascii/internal/scenario"
	"maps"
	"slices"
	"strconv"
	"strings"
)

type ControlNode struct {
	SelectableEntityID int
	TargetEntityID     int
	Next               *ControlNode
	Prev               *ControlNode
}

type World struct {
	ScenarioVariables        map[string]string
	Room                     string
	UserInputProfile         UserInputProfile
	InputProfiles            map[string]UserInputProfile
	InputProfileByRoom       map[string]string
	KeyDown                  string
	InfoboxContent           string
	ShouldQuit               bool
	ActiveControl            *ControlNode
	EditingControl           bool
	ControlOrder             map[string]*ControlNode
	EByPos                   map[component.Position]int
	Portals                  map[component.Position]component.Position
	Terminals                map[component.Position]string
	UILayout                 []string
	UIs                      []string
	UIContent                map[string][]string
	UIEmpty                  map[string][]string
	InfoboxRoomBaseContent   map[string]string
	HasChanged               bool
	IterationNr              int
	Entities                 []int
	Pos                      map[int]component.Position
	Layer                    map[int]component.Layer
	Ascii                    map[int]component.Ascii
	Impassable               map[int]component.Impassable
	Player                   map[int]component.Player
	Interactable             map[int]component.Interactable
	ControlNumber            map[int]component.ControlNumber
	ControlOptions           map[int]component.ControlOptions
	ControlLabels            map[int]component.ControlLabel
	Selectable               map[int]component.Selectable
	SelectableButtonVariable map[int]component.SelectableButtonVariable
}

func NewWorldEmpty() World {
	return World{
		ScenarioVariables:        make(map[string]string),
		InputProfiles:            make(map[string]UserInputProfile),
		InputProfileByRoom:       make(map[string]string),
		EByPos:                   make(map[component.Position]int),
		Portals:                  make(map[component.Position]component.Position),
		Terminals:                make(map[component.Position]string),
		UIContent:                make(map[string][]string),
		UIEmpty:                  make(map[string][]string),
		InfoboxRoomBaseContent:   make(map[string]string),
		Pos:                      make(map[int]component.Position),
		Layer:                    make(map[int]component.Layer),
		Ascii:                    make(map[int]component.Ascii),
		Impassable:               make(map[int]component.Impassable),
		Player:                   make(map[int]component.Player),
		Interactable:             make(map[int]component.Interactable),
		ControlNumber:            make(map[int]component.ControlNumber),
		ControlOptions:           make(map[int]component.ControlOptions),
		ControlLabels:            make(map[int]component.ControlLabel),
		Selectable:               make(map[int]component.Selectable),
		SelectableButtonVariable: make(map[int]component.SelectableButtonVariable),
		ControlOrder:             make(map[string]*ControlNode),
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
					if _, exists := entityIDsByName[entityName]; !exists {
						entityIDsByName[entityName] = entityID
					}
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
			if _, exists := entityIDsByName[entityName]; !exists {
				entityIDsByName[entityName] = entityID
			}
			if roomEntityIDs[roomName] == nil {
				roomEntityIDs[roomName] = make(map[rune]int)
			}
			roomEntityIDs[roomName][char] = entityID
		}
	}
	if err := collectControlLabels(&world); err != nil {
		return world, err
	}
	resolveControlLabelSources(&world, entityIDsByName)
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
				if _, ok := world.SelectableButtonVariable[entityID]; !ok {
					return world, fmt.Errorf("entity %q in room %q is not selectable", marker, roomName)
				}
			}
			ids = append(ids, entityID)
		}
		if len(ids) > 0 {
			selectables := maps.Clone(world.Selectable)
			for entityID, button := range world.SelectableButtonVariable {
				selectables[entityID] = component.Selectable{
					UnfocusedASCII: button.UnfocusedASCII,
					FocusedASCII:   button.FocusedASCII,
					SelectedASCII:  button.FocusedASCII,
				}
			}
			world.ControlOrder[roomName] = circularControlNodes(ids, selectables)
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
		var err error
		world, err = world.WithRoom(position.Room)
		if err != nil {
			return world, err
		}
	}
	return world, nil
}
func (w World) WithNextControl() (World, bool) {
	if w.ActiveControl == nil || w.ActiveControl.Next == nil {
		return w, false
	}
	next := w.Clone()
	next.ActiveControl = next.ActiveControl.Next
	return next, true
}

func (w World) WithPreviousControl() (World, bool) {
	if w.ActiveControl == nil || w.ActiveControl.Prev == nil {
		return w, false
	}
	next := w.Clone()
	next.ActiveControl = next.ActiveControl.Prev
	return next, true
}
func (w World) WithRoom(roomName string) (World, error) {
	next, ok := w.WithInputProfileForRoom(roomName)
	if !ok {
		return w, fmt.Errorf("input profile for room %q does not exist", roomName)
	}
	next.Room = roomName
	next.InfoboxContent = next.InfoboxRoomBaseContent[roomName]
	return next, nil
}

func (w World) WithInputProfileForRoom(roomName string) (World, bool) {
	profile, ok := w.InputProfileForRoom(roomName)
	if !ok {
		return w, false
	}
	next := w.Clone()
	next.UserInputProfile = profile
	next.EditingControl = false
	next.ActiveControl = next.ControlOrder[roomName]
	return next, true
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
	clone.ScenarioVariables = maps.Clone(w.ScenarioVariables)
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
	clone.ControlOptions = maps.Clone(w.ControlOptions)
	clone.ControlLabels = make(map[int]component.ControlLabel, len(w.ControlLabels))
	for entityID, label := range w.ControlLabels {
		label.EntityIDs = slices.Clone(label.EntityIDs)
		label.Sources = slices.Clone(label.Sources)
		label.SourceEntityIDs = slices.Clone(label.SourceEntityIDs)
		label.History = slices.Clone(label.History)
		clone.ControlLabels[entityID] = label
	}
	clone.Selectable = maps.Clone(w.Selectable)
	clone.SelectableButtonVariable = maps.Clone(w.SelectableButtonVariable)
	clone.ControlOrder = maps.Clone(w.ControlOrder)
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
	var controlOptions component.ControlOptions
	var hasControlOptions bool
	var controlLabel component.ControlLabel
	var hasControlLabel bool
	var selectable component.Selectable
	var hasSelectable bool
	var selectableButtonVariable component.SelectableButtonVariable
	var hasSelectableButtonVariable bool

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

		case component.NameControlOptions:
			if len(values) == 0 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			list := make([]rune, len(values))
			for i, value := range values {
				chars := []rune(value)
				if len(chars) != 1 {
					return fmt.Errorf("invalid values for component %q", name)
				}
				list[i] = chars[0]
			}
			controlOptions = component.ControlOptions{Current: list[0], Options: list}
			hasControlOptions = true

		case component.NameControlLabel:
			parsedLabel, err := parseControlLabel(values)
			if err != nil {
				return err
			}
			controlLabel = parsedLabel
			hasControlLabel = true

		case component.NameSelectableButtonVariable:
			if len(values) != 3 || values[2] == "" {
				return fmt.Errorf("invalid values for component %q", name)
			}
			unfocused := []rune(values[0])
			focused := []rune(values[1])
			if len(unfocused) != 1 || len(focused) != 1 {
				return fmt.Errorf("invalid values for component %q", name)
			}
			variable, value, ok := strings.Cut(values[2], "=")
			if !ok || variable == "" || value == "" || strings.Contains(value, "=") {
				return fmt.Errorf("invalid values for component %q", name)
			}
			selectableButtonVariable = component.SelectableButtonVariable{
				UnfocusedASCII: unfocused[0],
				FocusedASCII:   focused[0],
				VariableUpdate: values[2],
			}
			hasSelectableButtonVariable = true

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
	if hasControlOptions {
		w.ControlOptions[eID] = controlOptions
	}
	if hasControlLabel {
		w.ControlLabels[eID] = controlLabel
		if !hasASCII {
			w.Ascii[eID] = component.Ascii{Ascii: ' '}
		}
	}
	if hasSelectable {
		w.Selectable[eID] = selectable
	}
	if hasSelectableButtonVariable {
		w.SelectableButtonVariable[eID] = selectableButtonVariable
	}
	return nil
}

func parseControlLabel(values []string) (component.ControlLabel, error) {
	if len(values) != 2 {
		return component.ControlLabel{}, fmt.Errorf("invalid values for component %q", component.NameControlLabel)
	}
	width, err := strconv.Atoi(values[0])
	if err != nil || width <= 0 {
		return component.ControlLabel{}, fmt.Errorf("invalid values for component %q", component.NameControlLabel)
	}
	operation, sourcesText, ok := strings.Cut(values[1], ":")
	if !ok || operation == "" || !strings.HasPrefix(sourcesText, "[") || !strings.HasSuffix(sourcesText, "]") {
		return component.ControlLabel{}, fmt.Errorf("invalid values for component %q", component.NameControlLabel)
	}
	sources, err := splitControlLabelSources(strings.TrimSuffix(strings.TrimPrefix(sourcesText, "["), "]"))
	if err != nil {
		return component.ControlLabel{}, fmt.Errorf("invalid values for component %q: %w", component.NameControlLabel, err)
	}
	return component.ControlLabel{MaxLength: width, Operation: operation, Sources: sources}, nil
}

func splitControlLabelSources(text string) ([]string, error) {
	values := make([]string, 0)
	start := 0
	inQuotes := false
	escaped := false
	for i, char := range text {
		if inQuotes {
			if escaped {
				escaped = false
			} else if char == '\\' {
				escaped = true
			} else if char == '"' {
				inQuotes = false
			}
			continue
		}
		if char == '"' {
			inQuotes = true
		} else if char == ',' {
			values = append(values, strings.TrimSpace(text[start:i]))
			start = i + 1
		}
	}
	if inQuotes {
		return nil, fmt.Errorf("unterminated source string")
	}
	values = append(values, strings.TrimSpace(text[start:]))
	for i, value := range values {
		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
			unquoted, err := strconv.Unquote(value)
			if err != nil {
				return nil, err
			}
			values[i] = unquoted
		}
	}
	return values, nil
}

func collectControlLabels(w *World) error {
	groups := make(map[string][]int)
	for entityID, label := range w.ControlLabels {
		position := w.Pos[entityID]
		key := fmt.Sprintf("%s|%d|%s|%s", position.Room, label.MaxLength, label.Operation, strings.Join(label.Sources, "\x00"))
		groups[key] = append(groups[key], entityID)
	}
	for _, entityIDs := range groups {
		if len(entityIDs) == 0 {
			continue
		}
		label := w.ControlLabels[entityIDs[0]]
		slices.SortFunc(entityIDs, func(a, b int) int {
			pa, pb := w.Pos[a], w.Pos[b]
			if pa.Y != pb.Y {
				return pa.Y - pb.Y
			}
			return pa.X - pb.X
		})
		first := w.Pos[entityIDs[0]]
		minX, maxX := first.X, first.X
		minY, maxY := first.Y, first.Y
		for _, entityID := range entityIDs[1:] {
			position := w.Pos[entityID]
			minX, maxX = min(minX, position.X), max(maxX, position.X)
			minY, maxY = min(minY, position.Y), max(maxY, position.Y)
		}
		width := maxX - minX + 1
		height := maxY - minY + 1
		if len(entityIDs) != width*height {
			return fmt.Errorf("control label does not form a rectangle")
		}
		for i, entityID := range entityIDs {
			position := w.Pos[entityID]
			expected := component.Position{Room: first.Room, X: minX + i%width, Y: minY + i/width}
			if position != expected {
				return fmt.Errorf("control label entities do not form a rectangle")
			}
		}
		label.EntityIDs = slices.Clone(entityIDs)
		label.Width = width
		label.Height = height
		w.ControlLabels[entityIDs[0]] = label
		for _, entityID := range entityIDs[1:] {
			delete(w.ControlLabels, entityID)
		}
	}
	return nil
}

func resolveControlLabelSources(w *World, entityIDsByName map[string]int) {
	for entityID, label := range w.ControlLabels {
		label.SourceEntityIDs = make([]int, len(label.Sources))
		for i, source := range label.Sources {
			label.SourceEntityIDs[i] = -1
			if sourceEntityID, ok := entityIDsByName[source]; ok {
				label.SourceEntityIDs[i] = sourceEntityID
			}
		}
		w.ControlLabels[entityID] = label
	}
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
