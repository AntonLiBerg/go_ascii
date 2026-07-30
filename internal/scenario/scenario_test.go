package scenario

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAsciiMapFromMapFileContent(t *testing.T) {

	asciiMap := GetAsciiMap("===MAP\nab\ncd\n===ENTITY\nfirst\n- ascii:a")

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

	asciiMap := GetAsciiMap("å.\n#o")

	if got := asciiMap[[2]int{0, 0}]; got != 'å' {
		t.Fatalf("expected coordinate 0,0 to be 'å', got %q", got)
	}
	if got := asciiMap[[2]int{1, 1}]; got != 'o' {
		t.Fatalf("expected coordinate 1,1 to be 'o', got %q", got)
	}
}

func TestGetAsciiMapHandlesWindowsLineEndings(t *testing.T) {
	asciiMap := GetAsciiMap("===MAP\r\nab\r\ncd\r\n===ENTITY\r\nfirst\r\n- ascii:a")

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
	mapFile := "===layer[0]\n#.\n.#\n===layer[1]\n p\n"
	contentFile := "====ENTITY\n.:floor\n- pos\n- ascii:.\np:player\n- pos\n- ascii:o\n- player\n#:wall\n- pos\n- ascii=#\n- impassable\n====USERINPUTPROFILE\nquitgame=q\nmoveleft:a\n"

	if err := os.WriteFile(mapPath, []byte(mapFile), 0o644); err != nil {
		t.Fatalf("write temp map file: %v", err)
	}
	if err := os.WriteFile(contentPath, []byte(contentFile), 0o644); err != nil {
		t.Fatalf("write temp content file: %v", err)
	}

	asciiMap, entities, components, userInputProfileMap, err := GetScenarioFromFiles(mapPath, contentPath)
	if err != nil {
		t.Fatalf("GetScenarioFromFiles returned error: %v", err)
	}

	if len(asciiMap) != 2 {
		t.Fatalf("expected 2 layers, got %d", len(asciiMap))
	}
	if got := asciiMap[0][[2]int{0, 1}]; got != '.' {
		t.Fatalf("expected layer 0 coordinate 0,1 to be '.', got %q", got)
	}
	if got := asciiMap[1][[2]int{1, 0}]; got != 'p' {
		t.Fatalf("expected layer 1 coordinate 1,0 to be 'p', got %q", got)
	}
	if _, ok := asciiMap[1][[2]int{0, 0}]; ok {
		t.Fatal("expected spaces in higher layers to be transparent")
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

	if len(userInputProfileMap) != 2 {
		t.Fatalf("expected 2 user input profile entries, got %d", len(userInputProfileMap))
	}
	if got := userInputProfileMap["quitgame"]; got != "q" {
		t.Fatalf("expected quitgame button to be q, got %q", got)
	}
	if got := userInputProfileMap["moveleft"]; got != "a" {
		t.Fatalf("expected moveleft button to be a, got %q", got)
	}
}

func TestGetScenarioFromFilesParsesDoor(t *testing.T) {
	tempDir := t.TempDir()
	mapPath := filepath.Join(tempDir, "map.txt")
	contentPath := filepath.Join(tempDir, "content.txt")
	mapFile := "===layer[0]\n.\n===layer[1]\nd\n"
	contentFile := "====ENTITY\n.:floor\n- pos\n- ascii:.\nd:door\n- pos\n- ascii=D\n- impassable\n- door\n"

	if err := os.WriteFile(mapPath, []byte(mapFile), 0o644); err != nil {
		t.Fatalf("write temp map file: %v", err)
	}
	if err := os.WriteFile(contentPath, []byte(contentFile), 0o644); err != nil {
		t.Fatalf("write temp content file: %v", err)
	}

	_, entities, components, _, err := GetScenarioFromFiles(mapPath, contentPath)
	if err != nil {
		t.Fatalf("GetScenarioFromFiles returned error: %v", err)
	}

	if got := entities['d']; got != "door" {
		t.Fatalf("expected map key 'd' to be door, got %q", got)
	}
	assertComponentValues(t, components, "door", "pos")
	assertComponentValues(t, components, "door", "ascii", "D")
	assertComponentValues(t, components, "door", "impassable")
	assertComponentValues(t, components, "door", "door")
}

func TestGetScenarioFromFilesRejectsSkippedLayerNumber(t *testing.T) {
	_, err := GetLayeredAsciiMap("===layer[0]\n.\n===layer[2]\no")
	if err == nil {
		t.Fatal("expected skipped layer number error")
	}
}

func TestScenarioContentRejectsEntityWithoutKey(t *testing.T) {
	_, _, _, err := getEntitiesAndUserInputProfile("===ENTITY\nfloor\n- pos\n- ascii:.")
	if err == nil {
		t.Fatal("expected invalid entity header error")
	}
}

func TestScenarioContentParsesHyphenEntityKey(t *testing.T) {
	entities, components, _, err := getEntitiesAndUserInputProfile("===ENTITY\n-:wall\n- pos\n- ascii:#\n- impassable")
	if err != nil {
		t.Fatalf("getEntitiesAndUserInputProfile returned error: %v", err)
	}
	if got := entities['-']; got != "wall" {
		t.Fatalf("expected map key '-' to be wall, got %q", got)
	}
	assertComponentValues(t, components, "wall", "pos")
	assertComponentValues(t, components, "wall", "ascii", "#")
	assertComponentValues(t, components, "wall", "impassable")
}

func TestGetScenarioFromFilesReturnsError(t *testing.T) {
	tempDir := t.TempDir()
	asciiMap, entities, components, userInputProfileMap, err := GetScenarioFromFiles(
		filepath.Join(tempDir, "missing-map.txt"),
		filepath.Join(tempDir, "missing-content.txt"),
	)

	if err == nil {
		t.Fatal("expected error for missing map file")
	}
	if asciiMap != nil {
		t.Fatalf("expected nil asciiMap on error, got %v", asciiMap)
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
