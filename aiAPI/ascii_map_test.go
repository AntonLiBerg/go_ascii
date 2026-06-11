package aiapi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAsciiMapFromMapFileContent(t *testing.T) {
	api := AiAPI{}

	asciiMap := api.GetAsciiMap("===MAP\nab\ncd\n===ENTITY\na=first")

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
	api := AiAPI{}

	asciiMap := api.GetAsciiMap("å.\n#o")

	if got := asciiMap[[2]int{0, 0}]; got != 'å' {
		t.Fatalf("expected coordinate 0,0 to be 'å', got %q", got)
	}
	if got := asciiMap[[2]int{1, 1}]; got != 'o' {
		t.Fatalf("expected coordinate 1,1 to be 'o', got %q", got)
	}
}

func TestGetAsciiMapAndEntitiesFromFile(t *testing.T) {
	api := AiAPI{}
	tempDir := t.TempDir()
	mapPath := filepath.Join(tempDir, "map.txt")
	mapFile := "====MAP\n#.\no#\n====ENTITY\n.=floor\no=player\n#=wall\n"

	if err := os.WriteFile(mapPath, []byte(mapFile), 0o644); err != nil {
		t.Fatalf("write temp map file: %v", err)
	}

	asciiMap, entities := api.GetAsciiMapAndEntitiesFromFile(mapPath)

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
	if got := entities["floor"]; got != '.' {
		t.Fatalf("expected floor rune '.', got %q", got)
	}
	if got := entities["player"]; got != 'o' {
		t.Fatalf("expected player rune 'o', got %q", got)
	}
	if got := entities["wall"]; got != '#' {
		t.Fatalf("expected wall rune '#', got %q", got)
	}
}
