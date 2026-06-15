package aiapi

import (
	cmp "go_ascii/component"
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

func TestGetAsciiMapAndEntitiesFromFile(t *testing.T) {
	tempDir := t.TempDir()
	mapPath := filepath.Join(tempDir, "map.txt")
	mapFile := "====MAP\n#.\no#\n====ENTITY\nfloor\n- pos\n- ascii:.\n- tags: walkable, visible\nplayer\n- pos\n- ascii:o\nwall\n- pos\n- ascii=#\n"

	if err := os.WriteFile(mapPath, []byte(mapFile), 0o644); err != nil {
		t.Fatalf("write temp map file: %v", err)
	}

	asciiMap, entities, components, err := GetAsciiMapAndEntitiesFromFile(mapPath)
	if err != nil {
		t.Fatalf("GetAsciiMapAndEntitiesFromFile returned error: %v", err)
	}

	if len(asciiMap) != 4 {
		t.Fatalf("expected 4 map runes, got %d", len(asciiMap))
	}
	if got := asciiMap[[2]int{0, 1}]; got != 'o' {
		t.Fatalf("expected coordinate 0,1 to be 'o', got %q", got)
	}
	if got := asciiMap[[2]int{1, 0}]; got != '.' {
		t.Fatalf("expected coordinate 1,0 to be '.', got %q", got)
	}

	if len(entities) != 3 {
		t.Fatalf("expected 3 entities, got %d", len(entities))
	}
	if got := entities['.']; got != "floor" {
		t.Fatalf("expected rune '.' to be floor, got %q", got)
	}
	if got := entities['o']; got != "player" {
		t.Fatalf("expected rune 'o' to be player, got %q", got)
	}
	if got := entities['#']; got != "wall" {
		t.Fatalf("expected rune '#' to be wall, got %q", got)
	}

	if len(components) != 3 {
		t.Fatalf("expected 3 component entries, got %d", len(components))
	}
	assertComponentValues(t, components, "floor", cmp.C_POS)
	assertComponentValues(t, components, "floor", cmp.C_ASCII, ".")
	assertComponentValues(t, components, "floor", cmp.ComponentName("tags"), "walkable", "visible")
	assertComponentValues(t, components, "player", cmp.C_POS)
	assertComponentValues(t, components, "player", cmp.C_ASCII, "o")
	assertComponentValues(t, components, "wall", cmp.C_POS)
	assertComponentValues(t, components, "wall", cmp.C_ASCII, "#")
}

func TestGetAsciiMapAndEntitiesFromFileReturnsError(t *testing.T) {
	asciiMap, entities, components, err := GetAsciiMapAndEntitiesFromFile(filepath.Join(t.TempDir(), "missing.txt"))

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
}

func assertComponentValues(t *testing.T, components map[string]map[cmp.ComponentName][]string, entity string, component cmp.ComponentName, want ...string) {
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
