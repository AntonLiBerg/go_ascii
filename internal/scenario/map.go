package scenario

import (
	component "go_ascii/internal"
)

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
