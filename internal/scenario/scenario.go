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
	SectionNameDivider      string = "="
)

type Map struct {
	Rooms         map[string]map[[2]int]rune
	Ground        map[string]string
	InputProfiles map[string]string
	Portals       map[component.Position]component.Position
	Terminals     map[component.Position]string
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

func GetScenarioFromFiles(mapFilePath string, contentFilePath string) (Map, map[rune]string, map[string]map[string][]string, map[string]map[string]string, error) {
	mapContent, err := os.ReadFile(mapFilePath)
	if err != nil {
		return Map{}, nil, nil, nil, err
	}
	content, err := os.ReadFile(contentFilePath)
	if err != nil {
		return Map{}, nil, nil, nil, err
	}

	asciiMap, err := GetRoomMap(string(mapContent))
	if err != nil {
		return Map{}, nil, nil, nil, err
	}
	entities, components, inputProfiles, err := getEntitiesAndInputProfiles(string(content))
	if err != nil {
		return Map{}, nil, nil, nil, err
	}
	return asciiMap, entities, components, inputProfiles, nil
}

func GetRoomMap(mapText string) (Map, error) {
	asciiMap := Map{
		Rooms:         make(map[string]map[[2]int]rune),
		Ground:        make(map[string]string),
		InputProfiles: make(map[string]string),
		Portals:       make(map[component.Position]component.Position),
		Terminals:     make(map[component.Position]string),
	}
	mapText = strings.ReplaceAll(mapText, "\r\n", "\n")
	mapText = strings.ReplaceAll(mapText, "\r", "\n")
	lines := strings.Split(strings.TrimRight(mapText, "\n"), "\n")
	currentRoom := ""
	roomLines := []string(nil)
	inFeatures := false
	portalFeatures := [][4]string{}
	terminalFeatures := [][3]string{}

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
			currentRoom = strings.TrimSpace(strings.TrimPrefix(trimmed, "==="))
			if currentRoom == "" {
				return Map{}, fmt.Errorf("room header on line %d has no name", lineNumber+1)
			}
			if _, exists := asciiMap.Rooms[currentRoom]; exists {
				return Map{}, fmt.Errorf("duplicate room %q", currentRoom)
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

	return asciiMap, nil
}

func getEntitiesAndInputProfiles(contentText string) (map[rune]string, map[string]map[string][]string, map[string]map[string]string, error) {
	entities := make(map[rune]string)
	components := make(map[string]map[string][]string)
	inputProfiles := make(map[string]map[string]string)
	text := strings.ReplaceAll(contentText, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")

	entityStart := -1
	inputProfileStart := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		trimmed = strings.TrimLeft(trimmed, SectionNameDivider)
		switch trimmed {
		case SectionNameEntity:
			if entityStart == -1 {
				entityStart = i + 1
			}
		case SectionNameInputProfile:
			if inputProfileStart == -1 {
				inputProfileStart = i + 1
			}
		}
	}

	currentEntity := ""
	definedEntities := make(map[string]struct{})
	if entityStart != -1 {
		entityEnd := len(lines)
		for i := entityStart; i < len(lines); i++ {
			trimmed := strings.TrimSpace(lines[i])
			trimmed = strings.TrimLeft(trimmed, SectionNameDivider)
			if trimmed == SectionNameMap || trimmed == SectionNameEntity || trimmed == SectionNameInputProfile {
				entityEnd = i
				break
			}
		}

		for _, line := range lines[entityStart:entityEnd] {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
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
					components[currentEntity][componentName] = values
				}
				continue
			}

			keyText, name, ok := strings.Cut(line, ":")
			key := []rune(strings.TrimSpace(keyText))
			name = strings.TrimSpace(name)
			if !ok || len(key) != 1 || name == "" {
				return nil, nil, nil, fmt.Errorf("invalid entity header %q: expected key:name", line)
			}
			if existingName, exists := entities[key[0]]; exists {
				return nil, nil, nil, fmt.Errorf("duplicate entity key %q for %q and %q", key[0], existingName, name)
			}
			if _, exists := definedEntities[name]; exists {
				return nil, nil, nil, fmt.Errorf("duplicate entity name %q", name)
			}

			entities[key[0]] = name
			definedEntities[name] = struct{}{}
			components[name] = make(map[string][]string)
			currentEntity = name
		}
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
					return nil, nil, nil, fmt.Errorf("duplicate input profile %q", line)
				}
				inputProfiles[line] = make(map[string]string)
				currentProfile = line
				continue
			}
			if currentProfile == "" {
				return nil, nil, nil, fmt.Errorf("input binding %q has no profile", line)
			}

			binding := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			separator := strings.IndexAny(binding, ":=")
			if separator == -1 {
				return nil, nil, nil, fmt.Errorf("invalid input binding %q", line)
			}

			action := strings.TrimSpace(binding[:separator])
			button := strings.TrimSpace(binding[separator+1:])
			if action == "" || button == "" {
				return nil, nil, nil, fmt.Errorf("invalid input binding %q", line)
			}
			inputProfiles[currentProfile][action] = button
		}
	}

	return entities, components, inputProfiles, nil
}
