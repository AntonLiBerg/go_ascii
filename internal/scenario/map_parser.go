package scenario

import (
	"fmt"
	component "go_ascii/internal"
	"go_ascii/internal/helpers"
	"slices"
	"strings"
)

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
	

}
func makeFeaturesMap(lines []string)map[string][]string{
	fMap := make(map[string][]string)
	for _,line := range lines{
		_,l,_ := strings.Cut(line,"- ")
		name,vals,_ := strings.Cut(l,":")
	}
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
func isValidMapFile(roomLines []string) error {
	if len(roomLines) == 0 {
		return fmt.Errorf("roomLines is empty!")
	}
	if !isSectionHeader(roomLines[0]){
		return fmt.Errorf("incorrect header!")
	}
	if helpers.Any(roomLines,func(s string)bool{return s == ""}){
		return fmt.Errorf("room contains empty lines!")
	}
	return nil
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

	// Finish the current room by trimming separator lines and converting its
	// ASCII rows into coordinate-indexed runes.
	// Parse room headers and map rows first, collecting features that require
	// cross-room marker resolution for a later pass.
	for lineNumber, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "===") {
			if currentRoom != "" {
				if err := isValidMapFile(); err != nil {
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
	// Every room needs enough information to build its ground layer and select
	// an input profile before the world can be created.
	for roomName := range asciiMap.Rooms {
		if asciiMap.Ground[roomName] == "" {
			return Map{}, fmt.Errorf("room %q has no ground feature", roomName)
		}
		if asciiMap.InputProfiles[roomName] == "" {
			return Map{}, fmt.Errorf("room %q has no inputprofile feature", roomName)
		}
	}

	// Deferred features refer to markers rather than coordinates. Resolve each
	// marker once parsing has made every room available.
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

	// Portals are bidirectional, so both endpoints must be unique and are stored
	// as a pair of reverse mappings.
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

	// Terminals map one source marker to a target room without requiring a
	// corresponding marker in that target room.
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
