package scenario

import (
	"fmt"
	component "go_ascii/internal"
	"os"
	"strings"
)

const (
	SectionNameEntity       string = "ENTITY"
	SectionNameInputProfile string = "INPUTPROFILE"
	SectionNameMap          string = "MAP"
	SectionNameUILayout     string = "LAYOUT"
	SectionNameDivider      string = "="
	inputProfileTypeName           = "profiletype"
	inputProfileTypePrefix         = "profiletype"
)

type Map struct {
	Rooms           map[string]map[[2]int]rune
	Ground          map[string]string
	InputProfiles   map[string]string
	Portals         map[component.Position]component.Position
	Terminals       map[component.Position]string
	SelectableOrder map[string][]rune
	RoomGroups      map[string][]string
	EntityGroups    map[string]map[rune]struct{}
	UILayout        []string
	UIContent       map[string][]string
}

func GetAsciiMap(mapText string) map[[2]int]rune {
	asciiMap := make(map[[2]int]rune)
	mapText = strings.ReplaceAll(mapText, "\r\n", "\n")
	mapText = strings.ReplaceAll(mapText, "\r", "\n")
	lines := strings.Split(mapText, "\n")

	mapStart := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		trimmed = strings.TrimLeft(trimmed, SectionNameDivider)
		if trimmed == SectionNameMap {
			mapStart = i + 1
			break
		}
	}

	mapEnd := len(lines)
	if mapStart == -1 {
		mapStart = 0
	} else {
		for i := mapStart; i < len(lines); i++ {
			trimmed := strings.TrimSpace(lines[i])
			trimmed = strings.TrimLeft(trimmed, SectionNameDivider)
			if trimmed == SectionNameMap || trimmed == SectionNameEntity || trimmed == SectionNameInputProfile {
				mapEnd = i
				break
			}
		}
	}

	mapSection := strings.Trim(strings.Join(lines[mapStart:mapEnd], "\n"), "\n")
	if mapSection == "" {
		return asciiMap
	}

	for y, line := range strings.Split(mapSection, "\n") {
		// Spaces are transparent so lower layers remain visible.
		for x, char := range []rune(line) {
			asciiMap[[2]int{x, y}] = char
		}
	}

	return asciiMap
}

func GetScenarioFromFiles(mapFilePath string, contentFilePath string, uiFilePaths ...string) (Map, map[rune]string, map[string]map[string][]string, map[string]map[string]string, []string, []string, error) {
	if len(uiFilePaths) > 1 {
		return Map{}, nil, nil, nil, nil, nil, fmt.Errorf("expected at most one UI file path")
	}
	mapContent, err := os.ReadFile(mapFilePath)
	if err != nil {
		return Map{}, nil, nil, nil, nil, nil, err
	}
	content, err := os.ReadFile(contentFilePath)
	if err != nil {
		return Map{}, nil, nil, nil, nil, nil, err
	}
	asciiMap, err := GetRoomMap(string(mapContent))
	if err != nil {
		return Map{}, nil, nil, nil, nil, nil, err
	}
	entities, components, inputProfiles, groups, err := getEntitiesAndInputProfiles(string(content))
	if err != nil {
		return Map{}, nil, nil, nil, nil, nil, err
	}
	terminalRooms := make(map[string]struct{})
	for _, roomName := range asciiMap.Terminals {
		terminalRooms[roomName] = struct{}{}
	}
	// Compose each room's allowed entity set from all named groups.
	asciiMap.EntityGroups = make(map[string]map[rune]struct{})
	for roomName, groupNames := range asciiMap.RoomGroups {
		if len(groupNames) == 0 {
			continue
		}
		allowed := make(map[rune]struct{})
		for _, groupName := range groupNames {
			group, ok := groups[groupName]
			if !ok {
				if groupName == "terminal" {
					if _, isTerminal := terminalRooms[roomName]; isTerminal {
						continue
					}
				}
				return Map{}, nil, nil, nil, nil, nil, fmt.Errorf("group %q for room %q does not exist", groupName, roomName)
			}
			for entity := range group {
				allowed[entity] = struct{}{}
			}
		}
		asciiMap.EntityGroups[roomName] = allowed
	}
	if len(uiFilePaths) == 0 {
		return asciiMap, entities, components, inputProfiles, nil, nil, nil
	}
	uiFileContent, err := os.ReadFile(uiFilePaths[0])
	if err != nil {
		return Map{}, nil, nil, nil, nil, nil, err
	}
	layout, uis, uiContent, err := getUiLayoutAndUIs(string(uiFileContent))
	if err != nil {
		return Map{}, nil, nil, nil, nil, nil, err
	}
	asciiMap.UILayout = layout
	asciiMap.UIContent = uiContent
	return asciiMap, entities, components, inputProfiles, layout, uis, nil
}

func getUiLayoutAndUIs(uiFileText string) ([]string, []string, map[string][]string, error) {
	type section struct {
		name  string
		lines []string
	}

	text := strings.ReplaceAll(uiFileText, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")
	sections := []section{}
	current := section{}
	for lineNumber, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, SectionNameDivider+SectionNameDivider+SectionNameDivider) {
			if current.name != "" {
				sections = append(sections, current)
			}
			name := strings.TrimSpace(strings.TrimLeft(trimmed, SectionNameDivider))
			if name == "" {
				return nil, nil, nil, fmt.Errorf("UI section on line %d has no name", lineNumber+1)
			}
			current = section{name: name}
			continue
		}
		if current.name == "" {
			if trimmed != "" {
				return nil, nil, nil, fmt.Errorf("UI content must start with a section")
			}
			continue
		}
		current.lines = append(current.lines, line)
	}
	if current.name != "" {
		sections = append(sections, current)
	}

	var layout []string
	var uis []string
	uiContent := make(map[string][]string)
	for _, section := range sections {
		if strings.EqualFold(section.name, SectionNameUILayout) {
			if layout != nil {
				return nil, nil, nil, fmt.Errorf("duplicate UI layout section")
			}
			for lineNumber, line := range section.lines {
				name := strings.TrimSpace(line)
				if name == "" {
					continue
				}
				if strings.HasPrefix(name, "-") {
					return nil, nil, nil, fmt.Errorf("invalid UI layout entry on line %d", lineNumber+1)
				}
				for _, existing := range layout {
					if existing == name {
						return nil, nil, nil, fmt.Errorf("duplicate UI layout entry %q", name)
					}
				}
				layout = append(layout, name)
			}
			continue
		}

		for _, existing := range uis {
			if existing == section.name {
				return nil, nil, nil, fmt.Errorf("duplicate UI section %q", section.name)
			}
		}
		uis = append(uis, section.name)
		uiLines := section.lines
		// UI metadata is not part of the rendered lines.
		for i, line := range uiLines {
			if strings.TrimSpace(line) == "features" {
				uiLines = uiLines[:i]
				break
			}
		}
		for len(uiLines) > 0 && uiLines[0] == "" {
			uiLines = uiLines[1:]
		}
		for len(uiLines) > 0 && uiLines[len(uiLines)-1] == "" {
			uiLines = uiLines[:len(uiLines)-1]
		}
		uiContent[section.name] = uiLines
	}
	if layout == nil {
		return nil, nil, nil, fmt.Errorf("UI file has no layout section")
	}
	return layout, uis, uiContent, nil
}

func normalizeInputProfileType(value string) string {
	switch value {
	case "none", "terminal", "control":
		return inputProfileTypePrefix + value
	default:
		return value
	}
}

func GetRoomMap(mapText string) (Map, error) {
	asciiMap := Map{
		Rooms:           make(map[string]map[[2]int]rune),
		Ground:          make(map[string]string),
		InputProfiles:   make(map[string]string),
		Portals:         make(map[component.Position]component.Position),
		Terminals:       make(map[component.Position]string),
		SelectableOrder: make(map[string][]rune),
		RoomGroups:      make(map[string][]string),
		UIContent:       make(map[string][]string),
	}
	mapText = strings.ReplaceAll(mapText, "\r\n", "\n")
	mapText = strings.ReplaceAll(mapText, "\r", "\n")
	lines := strings.Split(strings.TrimRight(mapText, "\n"), "\n")
	currentRoom := ""
	roomLines := []string(nil)
	inFeatures := false
	portalFeatures := [][4]string{}
	terminalFeatures := [][3]string{}
	selectableOrderFeatures := make(map[string][]rune)

	storeRoom := func() error {
		for len(roomLines) > 0 && roomLines[0] == "" {
			roomLines = roomLines[1:]
		}
		for len(roomLines) > 0 && roomLines[len(roomLines)-1] == "" {
			roomLines = roomLines[:len(roomLines)-1]
		}
		if len(roomLines) == 0 {
			return fmt.Errorf("room %q has no map", currentRoom)
		}

		room := make(map[[2]int]rune)
		for y, line := range roomLines {
			for x, char := range []rune(line) {
				room[[2]int{x, y}] = char
			}
		}
		asciiMap.Rooms[currentRoom] = room
		return nil
	}

	for lineNumber, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "===") {
			if currentRoom != "" {
				if err := storeRoom(); err != nil {
					return Map{}, err
				}
			}
			roomHeader := strings.TrimSpace(strings.TrimPrefix(trimmed, "==="))
			roomName, groupText, hasGroups := strings.Cut(roomHeader, ":")
			currentRoom = strings.TrimSpace(roomName)
			if currentRoom == "" {
				return Map{}, fmt.Errorf("room header on line %d has no name", lineNumber+1)
			}
			if _, exists := asciiMap.Rooms[currentRoom]; exists {
				return Map{}, fmt.Errorf("duplicate room %q", currentRoom)
			}
			if hasGroups {
				for _, groupName := range strings.Split(groupText, ",") {
					groupName = strings.TrimSpace(groupName)
					if groupName == "" {
						return Map{}, fmt.Errorf("room %q has an empty entity group on line %d", currentRoom, lineNumber+1)
					}
					for _, existing := range asciiMap.RoomGroups[currentRoom] {
						if existing == groupName {
							return Map{}, fmt.Errorf("room %q repeats entity group %q", currentRoom, groupName)
						}
					}
					asciiMap.RoomGroups[currentRoom] = append(asciiMap.RoomGroups[currentRoom], groupName)
				}
			}
			roomLines = nil
			inFeatures = false
			continue
		}
		if currentRoom == "" {
			if trimmed == "" {
				continue
			}
			return Map{}, fmt.Errorf("map content must start with a room header")
		}
		if trimmed == "features" {
			inFeatures = true
			continue
		}
		if !inFeatures {
			roomLines = append(roomLines, line)
			continue
		}
		if trimmed == "" {
			continue
		}

		featureText, _, _ := strings.Cut(trimmed, "//")
		featureText = strings.TrimSpace(featureText)
		if !strings.HasPrefix(featureText, "- ") {
			return Map{}, fmt.Errorf("invalid feature on line %d", lineNumber+1)
		}
		name, value, ok := strings.Cut(strings.TrimSpace(strings.TrimPrefix(featureText, "-")), ":")
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		if !ok || value == "" {
			return Map{}, fmt.Errorf("invalid feature on line %d", lineNumber+1)
		}
		switch name {
		case "ground":
			if _, exists := asciiMap.Ground[currentRoom]; exists {
				return Map{}, fmt.Errorf("duplicate ground feature in room %q", currentRoom)
			}
			asciiMap.Ground[currentRoom] = value
		case "inputprofile":
			if _, exists := asciiMap.InputProfiles[currentRoom]; exists {
				return Map{}, fmt.Errorf("duplicate inputprofile feature in room %q", currentRoom)
			}
			asciiMap.InputProfiles[currentRoom] = value
		case "portal":
			values := strings.Split(value, ",")
			if len(values) != 3 {
				return Map{}, fmt.Errorf("invalid portal on line %d", lineNumber+1)
			}
			for i := range values {
				values[i] = strings.TrimSpace(values[i])
			}
			if len([]rune(values[0])) != 1 || values[1] == "" || len([]rune(values[2])) != 1 {
				return Map{}, fmt.Errorf("invalid portal on line %d", lineNumber+1)
			}
			portalFeatures = append(portalFeatures, [4]string{currentRoom, values[0], values[1], values[2]})
		case "terminal":
			values := strings.Split(value, ",")
			if len(values) != 2 {
				return Map{}, fmt.Errorf("invalid terminal on line %d", lineNumber+1)
			}
			for i := range values {
				values[i] = strings.TrimSpace(values[i])
			}
			if len([]rune(values[0])) != 1 || values[1] == "" {
				return Map{}, fmt.Errorf("invalid terminal on line %d", lineNumber+1)
			}
			terminalFeatures = append(terminalFeatures, [3]string{currentRoom, values[0], values[1]})
		case "selectableorder":
			if _, exists := selectableOrderFeatures[currentRoom]; exists {
				return Map{}, fmt.Errorf("duplicate selectableorder feature in room %q", currentRoom)
			}
			values := strings.Split(value, ",")
			markers := make([]rune, 0, len(values))
			for _, markerText := range values {
				markerText = strings.TrimSpace(markerText)
				marker := []rune(markerText)
				if len(marker) != 1 {
					return Map{}, fmt.Errorf("invalid selectableorder on line %d", lineNumber+1)
				}
				for _, existing := range markers {
					if existing == marker[0] {
						return Map{}, fmt.Errorf("selectableorder repeats marker %q on line %d", marker[0], lineNumber+1)
					}
				}
				markers = append(markers, marker[0])
			}
			selectableOrderFeatures[currentRoom] = markers
		default:
			return Map{}, fmt.Errorf("unknown feature %q on line %d", name, lineNumber+1)
		}
	}

	if currentRoom == "" {
		return Map{}, fmt.Errorf("map has no rooms")
	}
	if err := storeRoom(); err != nil {
		return Map{}, err
	}
	for roomName := range asciiMap.Rooms {
		if asciiMap.Ground[roomName] == "" {
			return Map{}, fmt.Errorf("room %q has no ground feature", roomName)
		}
		if asciiMap.InputProfiles[roomName] == "" {
			return Map{}, fmt.Errorf("room %q has no inputprofile feature", roomName)
		}
	}

	findMarker := func(roomName string, marker rune) (component.Position, error) {
		room, ok := asciiMap.Rooms[roomName]
		if !ok {
			return component.Position{}, fmt.Errorf("unknown room %q", roomName)
		}
		position := component.Position{}
		found := false
		for xy, char := range room {
			if char != marker {
				continue
			}
			if found {
				return component.Position{}, fmt.Errorf("marker %q appears more than once in room %q", marker, roomName)
			}
			position = component.Position{Room: roomName, X: xy[0], Y: xy[1]}
			found = true
		}
		if !found {
			return component.Position{}, fmt.Errorf("marker %q not found in room %q", marker, roomName)
		}
		return position, nil
	}

	for _, feature := range portalFeatures {
		sourceMarker := []rune(feature[1])[0]
		targetMarker := []rune(feature[3])[0]
		source, err := findMarker(feature[0], sourceMarker)
		if err != nil {
			return Map{}, err
		}
		target, err := findMarker(feature[2], targetMarker)
		if err != nil {
			return Map{}, err
		}
		if _, exists := asciiMap.Portals[source]; exists {
			return Map{}, fmt.Errorf("portal at %+v is already connected", source)
		}
		if _, exists := asciiMap.Portals[target]; exists {
			return Map{}, fmt.Errorf("portal at %+v is already connected", target)
		}
		asciiMap.Portals[source] = target
		asciiMap.Portals[target] = source
	}

	for _, feature := range terminalFeatures {
		sourceMarker := []rune(feature[1])[0]
		source, err := findMarker(feature[0], sourceMarker)
		if err != nil {
			return Map{}, err
		}
		if _, exists := asciiMap.Rooms[feature[2]]; !exists {
			return Map{}, fmt.Errorf("terminal references unknown room %q", feature[2])
		}
		if _, exists := asciiMap.Terminals[source]; exists {
			return Map{}, fmt.Errorf("terminal at %+v is already connected", source)
		}
		asciiMap.Terminals[source] = feature[2]
	}

	for roomName, markers := range selectableOrderFeatures {
		asciiMap.SelectableOrder[roomName] = append([]rune(nil), markers...)
	}

	return asciiMap, nil
}

func getEntitiesAndInputProfiles(contentText string) (map[rune]string, map[string]map[string][]string, map[string]map[string]string, map[string]map[rune]struct{}, error) {
	entities := make(map[rune]string)
	components := make(map[string]map[string][]string)
	inputProfiles := make(map[string]map[string]string)
	groups := make(map[string]map[rune]struct{})
	text := strings.ReplaceAll(contentText, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")

	inputProfileStart := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		trimmed = strings.TrimLeft(trimmed, SectionNameDivider)
		switch trimmed {
		case SectionNameInputProfile:
			if inputProfileStart == -1 {
				inputProfileStart = i + 1
			}
		}
	}

	currentGroup := ""
	currentEntity := ""
	definedEntities := make(map[string]struct{})
	entityEnd := len(lines)
	if inputProfileStart != -1 {
		entityEnd = inputProfileStart - 1
	}
	for lineNumber, line := range lines[:entityEnd] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "===") {
			groupName := strings.TrimSpace(strings.TrimLeft(line, SectionNameDivider))
			if groupName == "" {
				return nil, nil, nil, nil, fmt.Errorf("group header on line %d has no name", lineNumber+1)
			}
			if _, exists := groups[groupName]; exists {
				return nil, nil, nil, nil, fmt.Errorf("duplicate entity group %q", groupName)
			}
			groups[groupName] = make(map[rune]struct{})
			currentGroup = groupName
			currentEntity = ""
			continue
		}
		if currentGroup == "" {
			return nil, nil, nil, nil, fmt.Errorf("entity %q has no group", line)
		}

		if strings.HasPrefix(line, "- ") {
			componentText := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			if componentText == "" {
				continue
			}

			componentName := componentText
			values := []string{}
			separator := strings.IndexAny(componentText, ":=")
			if separator != -1 {
				name := strings.TrimSpace(componentText[:separator])
				if name == "" {
					continue
				}

				componentName = name
				valueText := strings.TrimSpace(componentText[separator+1:])
				if componentName == component.NameASCII && valueText == "SPACE" {
					values = append(values, " ")
				} else if valueText != "" {
					for _, value := range strings.Split(valueText, ",") {
						value = strings.TrimSpace(value)
						if value != "" {
							values = append(values, value)
						}
					}
				}
			}

			if currentEntity != "" {
				if _, exists := components[currentEntity][componentName]; exists {
					return nil, nil, nil, nil, fmt.Errorf("duplicate component %q for entity %q", componentName, currentEntity)
				}
				components[currentEntity][componentName] = values
			}
			continue
		}

		keyText, name, ok := strings.Cut(line, ":")
		key := []rune(strings.TrimSpace(keyText))
		name = strings.TrimSpace(name)
		if !ok || len(key) != 1 || name == "" {
			return nil, nil, nil, nil, fmt.Errorf("invalid entity header %q: expected key:name", line)
		}
		if existingName, exists := entities[key[0]]; exists {
			return nil, nil, nil, nil, fmt.Errorf("duplicate entity key %q for %q and %q", key[0], existingName, name)
		}
		if _, exists := definedEntities[name]; exists {
			return nil, nil, nil, nil, fmt.Errorf("duplicate entity name %q", name)
		}

		entities[key[0]] = name
		definedEntities[name] = struct{}{}
		components[name] = make(map[string][]string)
		groups[currentGroup][key[0]] = struct{}{}
		currentEntity = name
	}

	if inputProfileStart != -1 {
		inputProfileEnd := len(lines)
		for i := inputProfileStart; i < len(lines); i++ {
			trimmed := strings.TrimSpace(lines[i])
			trimmed = strings.TrimLeft(trimmed, SectionNameDivider)
			if trimmed == SectionNameMap || trimmed == SectionNameEntity || trimmed == SectionNameInputProfile {
				inputProfileEnd = i
				break
			}
		}

		currentProfile := ""
		for _, line := range lines[inputProfileStart:inputProfileEnd] {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if !strings.HasPrefix(line, "- ") {
				if _, exists := inputProfiles[line]; exists {
					return nil, nil, nil, nil, fmt.Errorf("duplicate input profile %q", line)
				}
				inputProfiles[line] = make(map[string]string)
				currentProfile = line
				continue
			}
			if currentProfile == "" {
				return nil, nil, nil, nil, fmt.Errorf("input binding %q has no profile", line)
			}

			binding := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			separator := strings.IndexAny(binding, ":=")
			if separator == -1 {
				return nil, nil, nil, nil, fmt.Errorf("invalid input binding %q", line)
			}

			action := strings.TrimSpace(binding[:separator])
			button := strings.TrimSpace(binding[separator+1:])
			if action == "" || button == "" {
				return nil, nil, nil, nil, fmt.Errorf("invalid input binding %q", line)
			}
			if _, exists := inputProfiles[currentProfile][action]; exists {
				return nil, nil, nil, nil, fmt.Errorf("duplicate input binding %q for profile %q", action, currentProfile)
			}
			if action == inputProfileTypeName {
				button = normalizeInputProfileType(button)
			}
			inputProfiles[currentProfile][action] = button
		}
	}

	return entities, components, inputProfiles, groups, nil
}
