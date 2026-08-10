package scenario

import (
	"fmt"
	"slices"
	"strings"
)
func isAllS1InS2(slice1 []string,slice2 []string)bool{
	return allIsTrue(slice1,func(s string)bool{	
		return slices.Contains(slice2,s)
	})
}
func isUnique(slice []string)bool{
	seen := make(map[string]bool)
	for _,s := range slice {
		if _,exists := seen[s]; exists{
			return false
		}
		seen[s] = true
	}
	return true
}
func normalize(slice []string)[]string{
	nSlice := filter(slice,func(s string)bool{
		return strings.TrimSpace(s) != ""
	})
	return nSlice
}
func isUiFileValid(lines []string) error{
	if len(lines) <= 1{
		return fmt.Errorf("empty file or just the header!")
	}
	if strings.TrimSpace(lines[0]) != SectionDivider+SectionNameUILayout{
		return fmt.Errorf("no layout section at start of file!")
	}
	uiLayout := []string{}
	remaining := lines[1:]
	for _,line := range remaining{
		if isUISectionHeader(line){
			break
		}
		trimmedLine := strings.TrimSpace(line)
		uiLayout = append(uiLayout, trimmedLine)
	}	
	uiNonLayoutSections := lines[len(uiLayout)+1:]
	sectionHeaders := filter(uiNonLayoutSections,func(s string) bool{
		return isUISectionHeader(s)
	})
	sectionHeaders = transformEach(sectionHeaders,func(h string)string{
		return strings.TrimSpace(strings.TrimPrefix(h, SectionDivider))
	})

	if !slices.Contains(uiLayout,"room"){
		return fmt.Errorf("layout does not contain required keyword room")
	}
	if !isUnique(uiLayout){
		return fmt.Errorf("duplicate headers in layout!")
	}
	if !isUnique(sectionHeaders){
		return fmt.Errorf("duplicate headers in ui file!")
	}
	if !isAllS1InS2(sectionHeaders,uiLayout){
		return fmt.Errorf("some sections are not mentioned in the layout!")
	}
	if !isAllS1InS2(filter(uiLayout,func(s string)bool{return s != "room"}),sectionHeaders){
		return fmt.Errorf("some sections mentioned in layout does not exist in the ui file!")
	}
	return nil
}
func hasAtMostOneUIFilePath(paths []string) bool {
	return len(paths) <= 1
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

func isUISectionHeader(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, SectionDivider) && strings.Count(line,SectionNameDivider) == 3
}

func isUIContentAllowed(sectionName string, line string) bool {
	return sectionName != "" || line == ""
}

func isUILayoutSection(name string) bool {
	return strings.EqualFold(name, SectionNameUILayout)
}

func isValidUILayoutEntry(name string) bool {
	return !strings.HasPrefix(name, "-")
}

func isUniqueUIName(names []string, name string) bool {
	for _, existing := range names {
		if existing == name {
			return false
		}
	}
	return true
}
