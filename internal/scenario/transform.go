package scenario

import (
	"fmt"
	component "go_ascii/internal"
	"slices"
	"strings"
)

func normalizeInputProfileType(value string) string {
	switch value {
	case "none", "terminal", "control":
		return inputProfileTypePrefix + value
	default:
		return value
	}
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
		trimmed := strings.TrimLeft(strings.TrimSpace(line), SectionNameDivider)
		if trimmed == SectionNameInputProfile && inputProfileStart == -1 {
			inputProfileStart = i + 1
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
			trimmed := strings.TrimLeft(strings.TrimSpace(lines[i]), SectionNameDivider)
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

func asciiMapFromLines(lines []string) map[[2]int]rune {
	asciiMap := make(map[[2]int]rune)
	for y, line := range lines {
		x := 0
		for _, char := range line {
			asciiMap[[2]int{x, y}] = char
			x++
		}
	}
	return asciiMap
}

func terminalRoomsFromTerminals(terminals map[component.Position]string) map[string]struct{} {
	rooms := make(map[string]struct{}, len(terminals))
	for _, roomName := range terminals {
		rooms[roomName] = struct{}{}
	}
	return rooms
}

func entityGroupsFromRoomGroups(roomGroups map[string][]string, groups map[string]map[rune]struct{}) map[string]map[rune]struct{} {
	entityGroups := make(map[string]map[rune]struct{}, len(roomGroups))
	for roomName, groupNames := range roomGroups {
		if len(groupNames) == 0 {
			continue
		}
		entities := make(map[rune]struct{})
		for _, groupName := range groupNames {
			for entity := range groups[groupName] {
				entities[entity] = struct{}{}
			}
		}
		entityGroups[roomName] = entities
	}
	return entityGroups
}

func withEntityGroups(asciiMap Map, entityGroups map[string]map[rune]struct{}) Map {
	asciiMap.EntityGroups = entityGroups
	return asciiMap
}

func withUI(asciiMap Map, layout []string, content map[string][]string) Map {
	asciiMap.UILayout = layout
	asciiMap.UIContent = content
	return asciiMap
}

func findMissingRoomGroup(roomGroups map[string][]string, groups map[string]map[rune]struct{}, terminalRooms map[string]struct{}) (string, string, bool) {
	for roomName, groupNames := range roomGroups {
		for _, groupName := range groupNames {
			if _, exists := groups[groupName]; exists {
				continue
			}
			if groupName == "terminal" {
				if _, isTerminal := terminalRooms[roomName]; isTerminal {
					continue
				}
			}
			return roomName, groupName, false
		}
	}
	return "", "", true
}

func findAsciiMapBounds(lines []string, sectionNames []string) (int, int) {
	lineCount := len(lines)
	start := 0
	found := false
	for i, name := range sectionNames {
		if name == SectionNameMap {
			start = i + 1
			found = true
			break
		}
	}
	end := lineCount
	if found {
		for i := start; i < lineCount; i++ {
			name := sectionNames[i]
			if name == SectionNameMap || name == SectionNameEntity || name == SectionNameInputProfile {
				end = i
				break
			}
		}
	}
	for start < end && lines[start] == "" {
		start++
	}
	for end > start && lines[end-1] == "" {
		end--
	}
	return start, end
}

func GetAsciiMap(mapText string) map[[2]int]rune {
	mapText = strings.ReplaceAll(mapText, "\r\n", "\n")
	mapText = strings.ReplaceAll(mapText, "\r", "\n")
	lines := strings.Split(mapText, "\n")
	sectionNames := make([]string, len(lines))
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		sectionNames[i] = strings.TrimLeft(trimmed, SectionNameDivider)
	}
	mapStart, mapEnd := findAsciiMapBounds(lines, sectionNames)
	return asciiMapFromLines(lines[mapStart:mapEnd])
}

func nGetRoomMap(mapText string) (Map, error){
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
	if err := isValidMapFile(lines); err != nil{
		return asciiMap, err
	}

	groupOnHeader := group(lines,func(s string)bool{return isSectionHeader(s)})
	for header,ls := range groupOnHeader{
		roomName,compositions := makeRoomHeaderParts(header)
		roomAscii := makeAsciiRoom(ls[:slices.Index(ls,"feature")])
		asciiMap.Rooms[roomName] = roomAscii
	}
	for header,ls := range groupOnHeader{
		features := strings.
	}


}
func makeFeaturesMap(lines []string)map[string][]string{
	fMap := make(map[string][]string)
	for _,line := range lines{
		_,l,_ := strings.Cut(line,"- ")
		name,vals,_ := strings.Cut(l,":")
		fMap[name] = strings.Split(vals, ",")
	}
	return fMap
}

func makeRoomHeaderParts(header string) (string,[]string){
	if !isSectionHeader(header){
		fmt.Errorf("not a header!")
	}
	_,afterDivider,_ :=strings.Cut(header,SectionDivider)
	name,appendedGroups,_ := strings.Cut(afterDivider,":")
	return name,strings.Split(appendedGroups,",")
}

func group(lines []string,f func(s string)bool)map[string][]string{
	currentHeader:= lines[0]
	lines = lines[1:]
	group := make(map[string][]string)
	for _,l := range lines{
		if f(l){
			currentHeader = l
			group[currentHeader] = []string{}
		}else{
			group[currentHeader] = append(group[currentHeader], l)
		}
	} 
	return group
}

func makeAsciiRoom(roomLines []string)map[[2]int]rune{
	room := make(map[[2]int]rune)
	for y, line := range roomLines {
		for x, char := range []rune(line) {
			room[[2]int{x, y}] = char
		}
	}
	return room
}


func getUiLayoutAndUIs(uiFileText string) ([]string, []string, map[string][]string, error) {
	text := strings.ReplaceAll(uiFileText, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")

	if err := isUiFileValid(lines); err != nil {
		return nil, nil, nil, err
	}
	uiLayout := getNextUiSection(lines[1:])
	uiSections := make(map[string][]string)
	remainingLines := lines[len(uiLayout)+1:]
	for len(remainingLines) > 0 {
		sName := strings.TrimSpace(strings.TrimPrefix(remainingLines[0], SectionDivider))
		uiSections[sName] = trimUIContent(getNextUiSection(remainingLines[1:]))
		remainingLines = remainingLines[len(uiSections[sName])+1:]
	}

	uis := make([]string, 0, len(uiLayout))
	for _, name := range uiLayout {
		if name != "room" {
			uis = append(uis, name)
		}
	}
	return uiLayout, uis, uiSections, nil
}

func trimUIContent(lines []string) []string {
	start := 0
	end := len(lines)
	for start < end && lines[start] == "" {
		start++
	}
	for end > start && lines[end-1] == "" {
		end--
	}
	return lines[start:end]
}

func getNextUiSection(lines []string) []string {
	section := []string{}
	for _, line := range lines {
		if isSectionHeader(line) {
			break
		}
		section = append(section, line)
	}
	return section
}
