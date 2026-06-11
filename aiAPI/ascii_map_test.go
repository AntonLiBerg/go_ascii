package aiapi

import "testing"

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
