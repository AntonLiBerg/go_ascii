package render

import (
	component "go_ascii/internal"
	"go_ascii/internal/world"
	"testing"
)

func TestPlayerCoversEntityAtSamePosition(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	position := [2]int{1, 1}
	if err := gameWorld.AddEntity(position, map[string][]string{
		"pos": {}, "ascii": {"D"}, "door": {},
	}); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}
	if err := gameWorld.AddEntity(position, map[string][]string{
		"pos": {}, "ascii": {"o"}, "player": {},
	}); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}

	pos := component.Position{X: 1, Y: 1}
	if !isCoveredByPlayer(gameWorld, 0, pos) {
		t.Fatal("expected door to be covered by player")
	}
	if isCoveredByPlayer(gameWorld, 1, pos) {
		t.Fatal("expected player to remain visible")
	}
}
