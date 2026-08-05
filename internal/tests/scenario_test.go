package tests

import (
	component "go_ascii/internal"
	"go_ascii/internal/scenario"
	"go_ascii/internal/world"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestGetAsciiMapFromMapFileContent(t *testing.T) {

	asciiMap := scenario.GetAsciiMap("===MAP\nab\ncd\n===ENTITY\nfirst\n- ascii:a")

	if len(asciiMap) != 4 {
		t.Fatalf("expected 4 map runes, got %d", len(asciiMap))
	}
	if got := asciiMap[[2]int{0, 0}]; got != 'a' {
		t.Fatalf("expected coordinate 0,0 to be 'a', got %q", got)
	}
	if got := asciiMap[[2]int{1, 0}]; got != 'b' {
		t.Fatalf("expected coordinate 1,0 to be 'b', got %q", got)
	}
	if got := asciiMap[[2]int{0, 1}]; got != 'c' {
		t.Fatalf("expected coordinate 0,1 to be 'c', got %q", got)
	}
	if got := asciiMap[[2]int{1, 1}]; got != 'd' {
		t.Fatalf("expected coordinate 1,1 to be 'd', got %q", got)
	}
}

func TestGetAsciiMapFromRawMapText(t *testing.T) {

	asciiMap := scenario.GetAsciiMap("å.\n#o")

	if got := asciiMap[[2]int{0, 0}]; got != 'å' {
		t.Fatalf("expected coordinate 0,0 to be 'å', got %q", got)
	}
	if got := asciiMap[[2]int{1, 1}]; got != 'o' {
		t.Fatalf("expected coordinate 1,1 to be 'o', got %q", got)
	}
}

func TestGetAsciiMapHandlesWindowsLineEndings(t *testing.T) {
	asciiMap := scenario.GetAsciiMap("===MAP\r\nab\r\ncd\r\n===ENTITY\r\nfirst\r\n- ascii:a")

	if got := asciiMap[[2]int{0, 0}]; got != 'a' {
		t.Fatalf("expected coordinate 0,0 to be 'a', got %q", got)
	}
	if got := asciiMap[[2]int{1, 1}]; got != 'd' {
		t.Fatalf("expected coordinate 1,1 to be 'd', got %q", got)
	}
}

func TestGetScenarioFromFiles(t *testing.T) {
	tempDir := t.TempDir()
	mapPath := filepath.Join(tempDir, "map.txt")
	contentPath := filepath.Join(tempDir, "content.txt")
	mapFile := "===ship\n#p\n.0\nfeatures\n- ground:floor\n- portal:0,engine room,1 // paired rooms\n- inputprofile:topdown\n===engine room\n\n1#\n .\nfeatures\n- ground:floor\n- inputprofile:topdown\n"
	contentFile := "====ENTITY\n.:floor\n- pos\n- ascii:.\np:player\n- pos\n- ascii:o\n- player\n#:wall\n- pos\n- ascii=#\n- impassable\n====INPUTPROFILE\ntopdown\n- profiletype=none\n- quitgame=q\n- moveleft:a\n"

	if err := os.WriteFile(mapPath, []byte(mapFile), 0o644); err != nil {
		t.Fatalf("write temp map file: %v", err)
	}
	if err := os.WriteFile(contentPath, []byte(contentFile), 0o644); err != nil {
		t.Fatalf("write temp content file: %v", err)
	}

	asciiMap, entities, components, userInputProfileMap, _, _, err := scenario.GetScenarioFromFiles(mapPath, contentPath)
	if err != nil {
		t.Fatalf("GetScenarioFromFiles returned error: %v", err)
	}

	if len(asciiMap.Rooms) != 2 {
		t.Fatalf("expected 2 rooms, got %d", len(asciiMap.Rooms))
	}
	if got := asciiMap.Rooms["ship"][[2]int{0, 1}]; got != '.' {
		t.Fatalf("expected ship coordinate 0,1 to be '.', got %q", got)
	}
	if got := asciiMap.Rooms["engine room"][[2]int{0, 0}]; got != '1' {
		t.Fatalf("expected leading blank separator to be ignored, got %q", got)
	}
	if got := asciiMap.Rooms["engine room"][[2]int{0, 1}]; got != ' ' {
		t.Fatalf("expected room spaces to be retained for ground, got %q", got)
	}
	if asciiMap.Ground["ship"] != "floor" || asciiMap.Ground["engine room"] != "floor" {
		t.Fatalf("expected floor ground in both rooms, got %v", asciiMap.Ground)
	}
	if asciiMap.InputProfiles["ship"] != "topdown" || asciiMap.InputProfiles["engine room"] != "topdown" {
		t.Fatalf("expected topdown profile in both rooms, got %v", asciiMap.InputProfiles)
	}
	shipPortal := component.Position{Room: "ship", X: 1, Y: 1}
	enginePortal := component.Position{Room: "engine room", X: 0, Y: 0}
	if got := asciiMap.Portals[shipPortal]; got != enginePortal {
		t.Fatalf("expected ship portal to lead to %+v, got %+v", enginePortal, got)
	}
	if got := asciiMap.Portals[enginePortal]; got != shipPortal {
		t.Fatalf("expected engine portal to lead back to %+v, got %+v", shipPortal, got)
	}

	if len(entities) != 3 {
		t.Fatalf("expected 3 entities, got %d", len(entities))
	}
	if got := entities['.']; got != "floor" {
		t.Fatalf("expected rune '.' to be floor, got %q", got)
	}
	if got := entities['p']; got != "player" {
		t.Fatalf("expected map key 'p' to be player, got %q", got)
	}
	if _, exists := entities['o']; exists {
		t.Fatal("expected rendered ASCII 'o' not to become a map key")
	}
	if got := entities['#']; got != "wall" {
		t.Fatalf("expected rune '#' to be wall, got %q", got)
	}

	if len(components) != 3 {
		t.Fatalf("expected 3 component entries, got %d", len(components))
	}
	assertComponentValues(t, components, "floor", "pos")
	assertComponentValues(t, components, "floor", "ascii", ".")
	assertComponentValues(t, components, "player", "pos")
	assertComponentValues(t, components, "player", "ascii", "o")
	assertComponentValues(t, components, "player", "player")
	assertComponentValues(t, components, "wall", "pos")
	assertComponentValues(t, components, "wall", "ascii", "#")
	assertComponentValues(t, components, "wall", "impassable")

	if len(userInputProfileMap) != 1 {
		t.Fatalf("expected 1 input profile, got %d", len(userInputProfileMap))
	}
	if got := userInputProfileMap["topdown"]["quitgame"]; got != "q" {
		t.Fatalf("expected quitgame button to be q, got %q", got)
	}
	if got := userInputProfileMap["topdown"]["moveleft"]; got != "a" {
		t.Fatalf("expected moveleft button to be a, got %q", got)
	}
	if got := userInputProfileMap["topdown"]["profiletype"]; got != world.ProfileTypeNone {
		t.Fatalf("expected profile type %q, got %q", world.ProfileTypeNone, got)
	}
}

func TestGetScenarioFromFilesParsesDoor(t *testing.T) {
	tempDir := t.TempDir()
	mapPath := filepath.Join(tempDir, "map.txt")
	contentPath := filepath.Join(tempDir, "content.txt")
	mapFile := "===room\n.d\nfeatures\n- ground:floor\n- inputprofile:topdown\n"
	contentFile := "====ENTITY\n.:floor\n- pos\n- ascii:.\nd:door\n- pos\n- ascii=D\n- impassable\n- interactable:door\n====INPUTPROFILE\ntopdown\n- interact=e\n"

	if err := os.WriteFile(mapPath, []byte(mapFile), 0o644); err != nil {
		t.Fatalf("write temp map file: %v", err)
	}
	if err := os.WriteFile(contentPath, []byte(contentFile), 0o644); err != nil {
		t.Fatalf("write temp content file: %v", err)
	}

	_, entities, components, _, _, _, err := scenario.GetScenarioFromFiles(mapPath, contentPath)
	if err != nil {
		t.Fatalf("GetScenarioFromFiles returned error: %v", err)
	}

	if got := entities['d']; got != "door" {
		t.Fatalf("expected map key 'd' to be door, got %q", got)
	}
	assertComponentValues(t, components, "door", "pos")
	assertComponentValues(t, components, "door", "ascii", "D")
	assertComponentValues(t, components, "door", "impassable")
	assertComponentValues(t, components, "door", "interactable", "door")
}

func TestGetRoomMapRejectsMissingPortalMarker(t *testing.T) {
	_, err := scenario.GetRoomMap("===first\n0\nfeatures\n- ground:floor\n- portal:0,second,1\n- inputprofile:topdown\n===second\n.\nfeatures\n- ground:floor\n- inputprofile:topdown")
	if err == nil {
		t.Fatal("expected missing portal marker error")
	}
}

func TestGetRoomMapRejectsDuplicateRoom(t *testing.T) {
	_, err := scenario.GetRoomMap("===room\n.\nfeatures\n- ground:floor\n- inputprofile:topdown\n===room\n.")
	if err == nil {
		t.Fatal("expected duplicate room error")
	}
}

func TestScenarioContentRejectsEntityWithoutKey(t *testing.T) {
	_, _, _, err := getScenarioContent(t, "===ENTITY\nfloor\n- pos\n- ascii:.")
	if err == nil {
		t.Fatal("expected invalid entity header error")
	}
}

func TestScenarioContentParsesHyphenEntityKey(t *testing.T) {
	entities, components, _, err := getScenarioContent(t, "===ENTITY\n-:wall\n- pos\n- ascii:#\n- impassable")
	if err != nil {
		t.Fatalf("getEntitiesAndInputProfiles returned error: %v", err)
	}
	if got := entities['-']; got != "wall" {
		t.Fatalf("expected map key '-' to be wall, got %q", got)
	}
	assertComponentValues(t, components, "wall", "pos")
	assertComponentValues(t, components, "wall", "ascii", "#")
	assertComponentValues(t, components, "wall", "impassable")
}

func TestScenarioContentParsesSpaceASCII(t *testing.T) {
	_, components, _, err := getScenarioContent(t, "===ENTITY\nv:void\n- pos\n- ascii:SPACE")
	if err != nil {
		t.Fatalf("getEntitiesAndInputProfiles returned error: %v", err)
	}
	assertComponentValues(t, components, "void", "ascii", " ")
}

func TestGetRoomMapParsesTerminal(t *testing.T) {
	asciiMap, err := scenario.GetRoomMap("===comms\nT\nfeatures\n- ground:floor\n- inputprofile:topdown\n- terminal:T,scan\n===scan\n+\nfeatures\n- ground:void\n- inputprofile:scan")
	if err != nil {
		t.Fatalf("GetRoomMap returned error: %v", err)
	}
	source := component.Position{Room: "comms", X: 0, Y: 0}
	if got := asciiMap.Terminals[source]; got != "scan" {
		t.Fatalf("expected terminal target scan, got %q", got)
	}
}

func TestGetRoomMapParsesEntityGroups(t *testing.T) {
	asciiMap, err := scenario.GetRoomMap("===bridge: base, instruments\n.\nfeatures\n- ground:floor\n- inputprofile:topdown")
	if err != nil {
		t.Fatalf("GetRoomMap returned error: %v", err)
	}
	if got := asciiMap.RoomGroups["bridge"]; !slices.Equal(got, []string{"base", "instruments"}) {
		t.Fatalf("expected bridge groups, got %v", got)
	}
}

func TestGetScenarioFromFilesReturnsError(t *testing.T) {
	tempDir := t.TempDir()
	asciiMap, entities, components, userInputProfileMap, _, _, err := scenario.GetScenarioFromFiles(
		filepath.Join(tempDir, "missing-map.txt"),
		filepath.Join(tempDir, "missing-content.txt"),
	)

	if err == nil {
		t.Fatal("expected error for missing map file")
	}
	if len(asciiMap.Rooms) != 0 || len(asciiMap.Ground) != 0 || len(asciiMap.InputProfiles) != 0 || len(asciiMap.Portals) != 0 || len(asciiMap.Terminals) != 0 {
		t.Fatalf("expected empty asciiMap on error, got %v", asciiMap)
	}
	if entities != nil {
		t.Fatalf("expected nil entities on error, got %v", entities)
	}
	if components != nil {
		t.Fatalf("expected nil components on error, got %v", components)
	}
	if userInputProfileMap != nil {
		t.Fatalf("expected nil userInputProfileMap on error, got %v", userInputProfileMap)
	}
}

func getScenarioContent(t *testing.T, content string) (map[rune]string, map[string]map[string][]string, map[string]map[string]string, error) {
	t.Helper()
	tempDir := t.TempDir()
	mapPath := filepath.Join(tempDir, "map.txt")
	contentPath := filepath.Join(tempDir, "content.txt")
	mapFile := "===room\n.\nfeatures\n- ground:floor\n- inputprofile:topdown\n"
	if err := os.WriteFile(mapPath, []byte(mapFile), 0o644); err != nil {
		t.Fatalf("write temp map file: %v", err)
	}
	if err := os.WriteFile(contentPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp content file: %v", err)
	}

	_, entities, components, inputProfiles, _, _, err := scenario.GetScenarioFromFiles(mapPath, contentPath)
	return entities, components, inputProfiles, err
}

func assertComponentValues(t *testing.T, components map[string]map[string][]string, entity string, component string, want ...string) {
	t.Helper()

	componentsForEntity, ok := components[entity]
	if !ok {
		t.Fatalf("expected components for entity %q", entity)
	}

	got, ok := componentsForEntity[component]
	if !ok {
		t.Fatalf("expected component %q for entity %q", component, entity)
	}

	if len(got) != len(want) {
		t.Fatalf("expected %s.%s values %v, got %v", entity, component, want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %s.%s values %v, got %v", entity, component, want, got)
		}
	}
}
