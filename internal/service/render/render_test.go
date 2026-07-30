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
	if !isCovered(gameWorld, 0, pos) {
		t.Fatal("expected door to be covered by player")
	}
	if isCovered(gameWorld, 1, pos) {
		t.Fatal("expected player to remain visible")
	}
}

func TestHigherLayerCoversLowerLayer(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	position := [2]int{1, 1}
	if err := gameWorld.AddEntityAtLayer(position, 0, map[string][]string{
		"pos": {}, "ascii": {"."},
	}); err != nil {
		t.Fatalf("AddEntityAtLayer returned error: %v", err)
	}
	if err := gameWorld.AddEntityAtLayer(position, 1, map[string][]string{
		"pos": {}, "ascii": {"D"}, "door": {},
	}); err != nil {
		t.Fatalf("AddEntityAtLayer returned error: %v", err)
	}

	pos := component.Position{X: 1, Y: 1}
	if !isCovered(gameWorld, 0, pos) {
		t.Fatal("expected layer 0 entity to be covered")
	}
	if isCovered(gameWorld, 1, pos) {
		t.Fatal("expected layer 1 entity to remain visible")
	}
}
