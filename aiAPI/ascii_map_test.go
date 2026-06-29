package aiapi

import (
	cmp "go_ascii/component"
	wrld "go_ascii/world"
	"os"
	"path/filepath"
	"slices"
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

func TestGetAsciiMapAndEntitiesFromFile(t *testing.T) {
	tempDir := t.TempDir()
	mapPath := filepath.Join(tempDir, "map.txt")
	mapFile := "====MAP\n#.\no#\n====ENTITY\nfloor\n- pos\n- ascii:.\n- tags: walkable, visible\nplayer\n- pos\n- ascii:o\nwall\n- pos\n- ascii=#\n- impassable\n====USERINPUTPROFILE\nquitgame=q\nmoveleft:a\n"

	if err := os.WriteFile(mapPath, []byte(mapFile), 0o644); err != nil {
		t.Fatalf("write temp map file: %v", err)
	}

	asciiMap, entities, components, userInputProfileMap, err := GetAsciiMapAndEntitiesFromFile(mapPath)
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
	assertComponentValues(t, components, "floor", cmp.C_TAGS, "walkable", "visible")
	assertComponentValues(t, components, "player", cmp.C_POS)
	assertComponentValues(t, components, "player", cmp.C_ASCII, "o")
	assertComponentValues(t, components, "wall", cmp.C_POS)
	assertComponentValues(t, components, "wall", cmp.C_ASCII, "#")
	assertComponentValues(t, components, "wall", cmp.C_IMPASSABLE)

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

func TestGetAsciiMapAndEntitiesFromFileReturnsError(t *testing.T) {
	asciiMap, entities, components, userInputProfileMap, err := GetAsciiMapAndEntitiesFromFile(filepath.Join(t.TempDir(), "missing.txt"))

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

func TestGetNeighbors(t *testing.T) {
	world := wrld.NewWorldEmpty()
	pId := world.MakeNewEntityId()
	e00 := world.MakeNewEntityId()
	e01 := world.MakeNewEntityId()
	e02 := world.MakeNewEntityId()
	e12 := world.MakeNewEntityId()
	e22 := world.MakeNewEntityId()
	e21 := world.MakeNewEntityId()
	e20 := world.MakeNewEntityId()
	e10 := world.MakeNewEntityId()
	world = world.AddPosition(pId, cmp.Position{X: 1, Y: 1})
	tests := []struct {
		name   string
		want   []int
		filter []cmp.ComponentName
		world  wrld.World
	}{
		{
			"noNeighbors", []int{}, []cmp.ComponentName{},
			world,
		},
		{
			"oneNeighbor00", []int{e00}, []cmp.ComponentName{},
			world.
				AddPosition(e00, cmp.Position{X: 0, Y: 0}),
		},
		{
			"oneNeighbor01", []int{e01}, []cmp.ComponentName{},
			world.
				AddPosition(e01, cmp.Position{X: 0, Y: 1}),
		},
		{
			"oneNeighbor02", []int{e02}, []cmp.ComponentName{},
			world.
				AddPosition(e02, cmp.Position{X: 0, Y: 2}),
		},
		{
			"oneNeighbor12", []int{e12}, []cmp.ComponentName{},
			world.
				AddPosition(e12, cmp.Position{X: 1, Y: 2}),
		},
		{
			"oneNeighbor22", []int{e22}, []cmp.ComponentName{},
			world.
				AddPosition(e22, cmp.Position{X: 2, Y: 2}),
		},
		{
			"oneNeighbor21", []int{e21}, []cmp.ComponentName{},
			world.
				AddPosition(e21, cmp.Position{X: 2, Y: 1}),
		},
		{
			"oneNeighbor20", []int{e20}, []cmp.ComponentName{},
			world.
				AddPosition(e20, cmp.Position{X: 2, Y: 0}),
		},
		{
			"oneNeighbor10", []int{e10}, []cmp.ComponentName{},
			world.
				AddPosition(e10, cmp.Position{X: 1, Y: 0}),
		},
		{
			"oneNeighborAll", []int{e00, e01, e02, e12, e22, e21, e20, e10}, []cmp.ComponentName{},
			world.
				AddPosition(e00, cmp.Position{X: 0, Y: 0}).
				AddPosition(e01, cmp.Position{X: 0, Y: 1}).
				AddPosition(e02, cmp.Position{X: 0, Y: 2}).
				AddPosition(e12, cmp.Position{X: 1, Y: 2}).
				AddPosition(e22, cmp.Position{X: 2, Y: 2}).
				AddPosition(e21, cmp.Position{X: 2, Y: 1}).
				AddPosition(e20, cmp.Position{X: 2, Y: 0}).
				AddPosition(e10, cmp.Position{X: 1, Y: 0}),
		},
		{
			"oneNeighborFiltered", []int{}, []cmp.ComponentName{cmp.C_IMPASSABLE},
			world.
				AddPosition(pId, cmp.Position{X: 1, Y: 0}).
				AddImpassable(pId),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := GetNeighbors(tt.world, pId, tt.filter)
			if !slices.Equal(tt.want, actual) {
				t.Errorf("got %v, want %v", actual, tt.want)
			}
		})
	}
}
