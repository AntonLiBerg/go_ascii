package scenario

import (
	"fmt"
	component "go_ascii/internal"
	"os"
)

const (
	SectionNameEntity       string = "ENTITY"
	SectionNameInputProfile string = "INPUTPROFILE"
	SectionNameMap          string = "MAP"
	SectionNameUILayout     string = "layout"
	SectionNameDivider      string = "="
	SectionDivider          string = SectionNameDivider + SectionNameDivider + SectionNameDivider
	inputProfileTypeName    string = "profiletype"
	inputProfileTypePrefix  string = "profiletype"
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

func GetScenarioFromFiles(mapFilePath string, contentFilePath string, uiFilePaths ...string) (Map, map[rune]string, map[string]map[string][]string, map[string]map[string]string, []string, []string, error) {
	if !hasAtMostOneUIFilePath(uiFilePaths) {
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
	terminalRooms := terminalRoomsFromTerminals(asciiMap.Terminals)
	roomName, groupName, valid := findMissingRoomGroup(asciiMap.RoomGroups, groups, terminalRooms)
	if !valid {
		return Map{}, nil, nil, nil, nil, nil, fmt.Errorf("group %q for room %q does not exist", groupName, roomName)
	}
	asciiMap = withEntityGroups(asciiMap, entityGroupsFromRoomGroups(asciiMap.RoomGroups, groups))
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
	asciiMap = withUI(asciiMap, layout, uiContent)
	return asciiMap, entities, components, inputProfiles, layout, uis, nil
}
