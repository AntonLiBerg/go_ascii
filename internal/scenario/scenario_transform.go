package scenario

import (
	"fmt"
	component "go_ascii/internal"
	"go_ascii/internal/helpers"
	"slices"
	"strings"
	AI "go_ascii/AI_workspace"

)

func normalizeInputProfileType(value string) string {
	switch value {
	case "none", "terminal", "control":
		return inputProfileTypePrefix + value
	default:
		return value
	}
}
type FileContent struct{
	entities      map[rune]string
	components    map[string]map[string][]string
	inputProfiles map[string]map[string]string
	groups        map[string]map[rune]struct{}
}
func ParseContentFile(contentText string) (FileContent,error) {
	contentMap := FileContent {
		entities:      make(map[rune]string),
		components:    make(map[string]map[string][]string),
		inputProfiles: make(map[string]map[string]string),
		groups:        make(map[string]map[rune]struct{}),
	}

	contentText = strings.ReplaceAll(contentText, "\r\n", "\n")
	contentText = strings.ReplaceAll(contentText, "\r", "\n")
	lines := strings.Split(strings.TrimRight(contentText, "\n"), "\n")
	if err := AI.IsValidContentFile(lines); err != nil{
		return contentMap,err
	}
	groupByName := group(lines,func(s string)bool{return isSectionHeader(s)},func(s string)string{
		return strings.Split(s,SectionDivider)[0]
	})
	for name,ls := range groupByName {
		if name == SectionNameInputProfile{
			if nInpProfiles,err := AI.MakeInputProfiles(ls); err != nil{
				return contentMap,err
			}else{
				contentMap.inputProfiles = nInpProfiles
			}
			continue
		}
		if nEntities,err := AI.AppendEntities(ls,contentMap.entities); err != nil{
			return contentMap,err
		}else{
			contentMap.entities = nEntities
		}
		if nComps,err := AI.AppendComponents(contentMap,ls);err != nil{
			return contentMap,err
		}else{
			contentMap.components = nComps
		}

		if nGroups,err := AI.AppendGroups(contentMap,name,ls);err != nil{
			return contentMap,err
		}else{
			contentMap.groups = nGroups
		}
	}
	return contentMap,nil
}

func GetRoomMap(mapText string) (FileMap, error) {
	asciiMap := FileMap{
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
	if err := isValidMapFile(lines); err != nil {
		return asciiMap, err
	}

	groupOnName := group(lines, func(s string) bool { return isSectionHeader(s) }, func(s string) string {
		roomName, _ := makeRoomHeaderParts(s)
		return roomName
	})
	for _, l := range lines {
		if isSectionHeader(l) {
			roomName, groups := makeRoomHeaderParts(l)
			asciiMap.RoomGroups[roomName] = groups
		}
	}
	for roomName, ls := range groupOnName {
		roomAscii := makeAsciiRoom(ls[:slices.Index(ls, "features")])
		asciiMap.Rooms[roomName] = roomAscii
	}
	if asciiMap, err := updateWithFeatures(groupOnName, asciiMap); err != nil {
		return asciiMap, err
	}

	return asciiMap, nil
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

func withEntityGroups(asciiMap FileMap, entityGroups map[string]map[rune]struct{}) FileMap {
	asciiMap.EntityGroups = entityGroups
	return asciiMap
}

func withUI(asciiMap FileMap, layout []string, content map[string][]string) FileMap {
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

func updateWithFeatures(groupOnName map[string][]string, asciiMap FileMap) (FileMap, error) {
	for roomName, ls := range groupOnName {
		features := makeFeaturesMap(helpers.SubSliceAfter(ls, "features"))

		for f, val := range features {
			switch f {
			case "ground":
				asciiMap.Ground[roomName] = val[0]
				break
			case "inputprofile":
				asciiMap.InputProfiles[roomName] = val[0]
				break
			case "portal":
				posFrom, err := getPosPortal(asciiMap.Rooms[roomName], val[0])
				if err != nil {
					return asciiMap, fmt.Errorf("portal %q in room %q: %w", val[0], roomName, err)
				}
				posTo, err := getPosPortal(asciiMap.Rooms[val[1]], val[2])
				if err != nil {
					return asciiMap, fmt.Errorf("portal %q in room %q: %w", val[2], val[1], err)
				}
				from := component.Position{Room: roomName, X: posFrom[0], Y: posFrom[1]}
				to := component.Position{Room: val[1], X: posTo[0], Y: posTo[1]}
				asciiMap.Portals[from] = to
				asciiMap.Portals[to] = from
			case "terminal":
				termPos, err := getPosPortal(asciiMap.Rooms[roomName], val[0])
				if err != nil {
					return asciiMap, fmt.Errorf("terminal %q in room %q: %w", val[0], roomName, err)
				}
				tPos := component.Position{Room: roomName, X: termPos[0], Y: termPos[1]}
				asciiMap.Terminals[tPos] = val[1]
				break
			case "selectableorder":
				markers := make([]rune, 0, len(val))
				for _, marker := range val {
					marker = strings.TrimSpace(marker)
					runes := []rune(marker)
					if len(runes) != 1 {
						return asciiMap, fmt.Errorf("invalid selectable marker %q in room %q: expected one character", marker, roomName)
					}
					markers = append(markers, runes[0])
				}
				asciiMap.SelectableOrder[roomName] = markers
			}
		}
	}
	return asciiMap, nil
}
func makeFeaturesMap(lines []string) map[string][]string {
	fMap := make(map[string][]string)
	for _, line := range lines {
		_, l, _ := strings.Cut(line, "- ")
		name, vals, _ := strings.Cut(l, ":")
		fMap[name] = strings.Split(vals, ",")
	}
	return fMap
}

func makeRoomHeaderParts(header string) (string, []string) {
	if !isSectionHeader(header) {
		fmt.Errorf("not a header!")
	}
	_, afterDivider, _ := strings.Cut(header, SectionDivider)
	name, appendedGroups, _ := strings.Cut(afterDivider, ":")
	return name, strings.Split(appendedGroups, ",")
}

func group(lines []string, f func(s string) bool, ft func(s string) string) map[string][]string {
	currentHeader := ""
	group := make(map[string][]string)
	for _, l := range lines {
		if f(l) {
			currentHeader = ft(l)
			group[currentHeader] = []string{}
		} else {
			group[currentHeader] = append(group[currentHeader], l)
		}
	}
	return group
}

func getPosPortal(room map[[2]int]rune, portal string) ([2]int, error) {
	marker := []rune(strings.TrimSpace(portal))
	if len(marker) == 0 {
		return [2]int{}, fmt.Errorf("portal marker is empty")
	}
	for position, value := range room {
		if value == marker[0] {
			return position, nil
		}
	}
	return [2]int{}, fmt.Errorf("marker %q not found", string(marker[0]))
}

func makeAsciiRoom(roomLines []string) map[[2]int]rune {
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
