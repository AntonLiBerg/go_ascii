package scenario

import (
	component "go_ascii/internal"
	"strings"
)

func filter[T any](slice []T,f func(T)bool) []T{
	nSlice := []T{}
	for _,s := range slice{
		if f(s){
			nSlice = append(nSlice, s)
		}
	}
	return nSlice
}
func anyIsTrue[T any](slice []T,f func(T)bool) bool{
	for _,s := range slice {
		if f(s){
			return true
		}
	}
	return false
}
func allIsTrue[T any](slice []T,f func(T)bool) bool{
	for _,s := range slice {
		if !f(s){
			return false
		}
	}
	return true
}
func transformEach[T any](values []T, transform func(T) T) []T {
	result := make([]T, 0, len(values))
	for _, value := range values {
		result = append(result, transform(value))
	}
	return result
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

func uiSectionNameFromHeader(header string) string {
	return strings.TrimSpace(strings.TrimLeft(header, SectionNameDivider))
}

func uiLayoutEntriesFromLines(lines []string) []uiLayoutEntry {
	entries := make([]uiLayoutEntry, 0, len(lines))
	for lineNumber, line := range lines {
		name := strings.TrimSpace(line)
		if name != "" {
			entries = append(entries, uiLayoutEntry{name: name, lineNumber: lineNumber})
		}
	}
	return entries
}

func uiContentFromLines(lines []string) []string {
	end := len(lines)
	for i, line := range lines {
		if strings.TrimSpace(line) == "features" {
			end = i
			break
		}
	}
	start := 0
	for start < end && lines[start] == "" {
		start++
	}
	for end > start && lines[end-1] == "" {
		end--
	}
	return lines[start:end]
}
